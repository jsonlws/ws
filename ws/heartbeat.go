package ws

import (
	"errors"
	"sync"
	"time"
)

//心跳处理程序 采用动态分组算法进行处理

/**
这里有两个指针：【当前时间】指针和【超时节点】指针
1、 【当前时间】指针指向的节点是  根据  curr_time % array_length 计算的，即 当前时间（秒时间戳） 对数组长度取模后的数组下标；
2、 【超时节点】指针指向的节点是  根据  (curr_time + 30) % array_length 计算的（假设心跳超时时间是 30秒，
即连续 30秒没有发心跳的客户端，可以判定为超时）。

正常发心跳的客户端会一直超出【当前时间】节点的30个位置。

【当前时间】指针会每秒钟向前走一个节点，节点指向的列表就是没有按时发心跳被扫描到的客户端。
*/

//定义桶的大小
const (
	bucketSize = 120
	outTime    = 30 //超时时间单位秒
)

//定义桶内容
type HeartBucket struct {
	HeartBeatTime int64 //最近一次心跳发送时间
	User          *User //ws中的用户连接信息
}

type HeartBucketLink struct {
	Data map[uint]*HeartBucket
	Lock *sync.RWMutex //定义读写锁
}

type Bucket struct {
	BucketLink       []*HeartBucketLink
	CurrentTimeIndex uint //当前时间的桶索引
	OutTimeIndex     uint //超时时间的桶索引
}

func NewHeartBeat() *Bucket {

	bucket := make([]*HeartBucketLink, bucketSize)

	for k, _ := range bucket {
		linkData := &HeartBucketLink{
			Data: make(map[uint]*HeartBucket),
			Lock: new(sync.RWMutex),
		}
		bucket[k] = linkData
	}

	//当前时间
	currentTime := time.Now().Unix()

	//获取当前时间的取模索引值
	currentTimeIndex := currentTime % bucketSize

	//超时时间取模值
	outTimeIndex := (currentTime + outTime) % bucketSize

	return &Bucket{
		BucketLink:       bucket,
		CurrentTimeIndex: uint(currentTimeIndex),
		OutTimeIndex:     uint(outTimeIndex),
	}
}

//客户端第一次链接到ws时
func (b *Bucket) FirstHeartHandler(clientId uint, user *User) uint {
	nowTime := time.Now().Unix()
	//计算出放置的桶索引节点
	index := (b.CurrentTimeIndex + outTime - 1) % bucketSize
	//因为map时非数据安全需要加锁处理
	b.BucketLink[index].Lock.Lock()
	defer b.BucketLink[index].Lock.Unlock()
	//写入对应桶的map中以客户端id为key
	b.BucketLink[index].Data[clientId] = &HeartBucket{
		HeartBeatTime: nowTime,
		User:          user,
	}

	return index
}

//之后正常心跳处理
func (b *Bucket) FutureHeartHandler(clientId uint, oldIndex uint, user *User) (uint, error) {

	if oldIndex >= bucketSize {
		return 0, errors.New("客户端携带原数据错误")
	}

	nowTime := time.Now().Unix()
	//1.先删除旧桶的数据

	//为了防止客户端进行心跳攻击这里需要判断如果用户心跳发送间隔时间时小于超时时间一半的时间则判定为异常

	if v, ok := b.BucketLink[oldIndex].Data[clientId]; ok {

		delete(b.BucketLink[oldIndex].Data, clientId)

		if (nowTime - v.HeartBeatTime) < int64(outTime/2) {
			return 0, errors.New("心跳发送频繁")
		}

	}

	//2.然后将数据放到新的桶中

	//新的索引节点
	bucketIndex := (b.CurrentTimeIndex + outTime - 1) % bucketSize

	b.BucketLink[bucketIndex].Lock.Lock()
	defer b.BucketLink[bucketIndex].Lock.Unlock()

	b.BucketLink[bucketIndex].Data[clientId] = &HeartBucket{
		HeartBeatTime: nowTime,
		User:          user,
	}

	return bucketIndex, nil
}

//拨动时针
func (b *Bucket) TurnClockwise(h *Hub) {
	for {
		//模拟时间钟每秒走一次
		time.Sleep(1 * time.Second)

		//超时时间索引扫描到时数据也是为超时链接
		outTimeLists1 := b.BucketLink[b.OutTimeIndex].Data
		for k, v := range outTimeLists1 {
			delete(outTimeLists1, k)
			h.unregister <- v.User
		}

		//当前时间索引扫描到的数据为超时链接
		outTimeLists2 := b.BucketLink[b.CurrentTimeIndex].Data
		for k, v := range outTimeLists2 {
			delete(outTimeLists2, k)
			h.unregister <- v.User
		}

		//每走一秒索引往前加1
		b.CurrentTimeIndex = (b.CurrentTimeIndex + 1) % bucketSize
		b.OutTimeIndex = (b.OutTimeIndex + 1) % bucketSize
	}
}
