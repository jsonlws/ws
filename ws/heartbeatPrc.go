package ws

import (
	"errors"
	"log"
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
	bucketSize = 120 //桶大小定义
	outTime    = 60  //超时时间单位秒
)

//心跳内容定义
type HeartBucket struct {
	ClientId      string //用户唯一标识id
	HeartBeatTime int64  //最近一次心跳时间
	User          *User  //用户ws中的连接信息
}

//心跳链表
type HeartBucketLink struct {
	Data map[string]HeartBucket
}

//桶管理器
type Bucket struct {
	BucketLink       []HeartBucketLink //桶列表
	CurrentTimeIndex uint              //当前时间的桶索引
	OutTimeIndex     uint              //超时时间的桶索引
}

func NewHeartBeat() *Bucket {
	//当前时间
	currentTime := time.Now().Unix()

	//获取当前时间的取模索引值
	currentTimeIndex := currentTime % bucketSize

	//超时时间取模值
	outTimeIndex := (currentTime + outTime) % bucketSize

	return &Bucket{
		BucketLink:       make([]HeartBucketLink, bucketSize),
		CurrentTimeIndex: uint(currentTimeIndex),
		OutTimeIndex:     uint(outTimeIndex),
	}
}

//客户端第一次链接到ws时map为不安全数据需要加锁处理
func (b *Bucket) FirstHeartHandler(clientId string, user *User, lock *sync.Mutex) uint {
	lock.Lock()
	defer lock.Unlock()

	nowTime := time.Now().Unix()

	bucketData := make(map[string]HeartBucket)

	bucketData[clientId] = HeartBucket{
		ClientId:      clientId,
		HeartBeatTime: nowTime,
		User:          user,
	}

	index := (b.CurrentTimeIndex + outTime - 1) % bucketSize

	//index := nowTime % bucketSize

	b.BucketLink[index] = HeartBucketLink{
		Data: bucketData,
	}

	return uint(index)
}

//之后正常心跳处理
func (b *Bucket) FutureHeartHandler(clientId string, oldIndex uint, user *User, lock *sync.Mutex) (uint, error) {
	lock.Lock()
	defer lock.Unlock()

	nowTime := time.Now().Unix()
	if oldIndex >= bucketSize {
		return 0, errors.New("客户端携带原数据错误")
	}
	//1.先删除旧桶的数据
	oldData := b.BucketLink[oldIndex]
	delete(oldData.Data, clientId)
	//2.然后将数据放到新的桶中
	bucketIndex := (b.CurrentTimeIndex + outTime - 1) % bucketSize

	//bucketIndex := nowTime % bucketSize

	bucketData := make(map[string]HeartBucket)

	bucketData[clientId] = HeartBucket{
		ClientId:      clientId,
		HeartBeatTime: nowTime,
		User:          user,
	}
	b.BucketLink[bucketIndex] = HeartBucketLink{
		Data: bucketData,
	}

	return uint(bucketIndex), nil
}

//拨动时针
func (b *Bucket) TurnClockwise(h *Hub) {
	for {
		//模拟时间钟每秒走一次
		time.Sleep(1 * time.Second)
		nowTime := time.Now().Unix()
		//超时时间索引扫描到时数据也是为超时链接
		outTimeLists1 := b.BucketLink[b.OutTimeIndex].Data

		for k, v := range outTimeLists1 {
			delete(outTimeLists1, k)
			if nowTime-v.HeartBeatTime > outTime {
				log.Println("#1用户", v.User.Uid, "超时未发送心跳")
				h.unregister <- v.User
			}
		}

		//当前时间索引扫描到的数据为超时链接
		outTimeLists2 := b.BucketLink[b.CurrentTimeIndex].Data
		for k, v := range outTimeLists2 {
			delete(outTimeLists2, k)
			if nowTime-v.HeartBeatTime > outTime {
				log.Println("#2用户", v.User.Uid, "超时未发送心跳")
				h.unregister <- v.User
			}
		}

		//每走一秒索引往前加1
		b.CurrentTimeIndex = (b.CurrentTimeIndex + 1) % bucketSize
		b.OutTimeIndex = (b.OutTimeIndex + 1) % bucketSize
	}
}
