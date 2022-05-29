package beanstalkd

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kr/beanstalk"
	"github.com/spf13/viper"
)

const (
	TubePayNotice        = "pay_notice"         //支付提醒通知队列
	TubePay              = "pay"                //支付队列
	TubeQueryPayOrder    = "query_pay_order"    //查询订单队列
	TubeQueryRefundOrder = "query_refund_order" //查询订单队列
)

//beanstalk消息投递所传参数结构体
type BeanstalkProducerParams struct {
	TubeName   string        //队列名称
	PutMsgBody []byte        //队列消息体
	Delay      time.Duration //延时执行时间
	Pri        uint32        //优先级
	Ttr        time.Duration //数据过期时间
}

func (b *BeanstalkProducerParams) Producer() (uint64, error) {
	c, err := beanstalk.Dial("tcp", viper.GetString("beanstalkd.addr"))
	if err != nil {
		return 0, err
	}
	defer c.Close()
	c.Tube.Name = b.TubeName
	c.TubeSet.Name[b.TubeName] = true
	jobId, err := c.Put(b.PutMsgBody, b.Pri, b.Delay, b.Ttr)
	return jobId, err
}

//beanstalk删除消息所传参数结构体
type BeanstalkDelJobParams struct {
	TubeName string //队列名称
	JobId    uint64 //任务id
}

func (b *BeanstalkDelJobParams) DelJob() error {
	c, err := beanstalk.Dial("tcp", viper.GetString("beanstalkd.addr"))
	if err != nil {
		return err
	}
	defer c.Close()
	c.Tube.Name = b.TubeName
	c.TubeSet.Name[b.TubeName] = true
	delErr := c.Delete(b.JobId)
	return delErr
}

//beanstalk消费消息所传参数结构体
type BeanstalkConsumerParams struct {
	TubeName     string                          //队列名称
	HandleMethod func(body []byte) (bool, error) //处理函数
	Limit        int                             //失败重试次数 0不做限制
	Frequency    time.Duration                   //从队列中取数据的频率
}

//若任务执行失败延时重试时间间隔
var TimeDelayIsSet map[int]time.Duration = map[int]time.Duration{
	1: 5 * time.Second,
	2: 15 * time.Second,
	3: 30 * time.Second,
	4: 60 * time.Second,
	5: 90 * time.Second,
}

func (b *BeanstalkConsumerParams) Consumer(ctx context.Context) {
	c, err := beanstalk.Dial("tcp", viper.GetString("beanstalkd.addr"))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.Tube.Name = b.TubeName
	c.TubeSet.Name[b.TubeName] = true
	b.handelConsumer(ctx, c)
}

func (b *BeanstalkConsumerParams) handelConsumer(ctx context.Context, c *beanstalk.Conn) {
	for {

		select {
		case <-ctx.Done():
			log.Printf("从beanstalk中读取数据到chan中的协程退出, because %v\n", ctx.Err())
			return
		default:

			//从队列中按设置取数据的频率取一次数据
			id, body, err := c.Reserve(b.Frequency)
			//说明此时队列中没有数据
			if id == 0 {
				continue
			}
			if err != nil {
				if !strings.Contains(err.Error(), "timeout") {
					log.Println(" [Consumer] [", c.Tube.Name, "] err:", err, " 任务id:", id)
				}
				continue
			}

			b.handelJob(c, body, id)

		}
	}
}

//处理任务数据
func (b *BeanstalkConsumerParams) handelJob(c *beanstalk.Conn, body []byte, id uint64) {
	//第一步校验数据格式是否符合规定
	sendBody := map[string]interface{}{}
	err := json.Unmarshal(body, &sendBody)
	if err != nil {
		c.Delete(id)
		log.Println(" [Consumer] [", c.Tube.Name, "] err:body参数必须为json格式", " 任务id:", id, "数据为:", string(body))
		return
	}

	//第二步获取当前job的统计信息
	jobInfo, _ := c.StatsJob(id)
	//若设置了失败重试次数当达到阈值就将该任务给删除，防止堆积
	jobRunNum, _ := strconv.Atoi(jobInfo["reserves"])
	if b.Limit > 0 && (jobRunNum-1) >= b.Limit {
		c.Delete(id)
		log.Println(" [Consumer] [", c.Tube.Name, "] 失败重试执行次数已经达到:", b.Limit, "已删除 任务id:", id, "数据为:", string(body))
		return
	}
	//调用对应函数处理
	ret, runerrs := b.HandleMethod(body)
	if ret {
		c.Delete(id)
		log.Println(" [Consumer] [", c.Tube.Name, "] 任务执行成功:", b.Limit, "任务id:", id, "数据为:", string(body))
	} else {
		//重新设置任务的延时
		if v, ok := TimeDelayIsSet[jobRunNum]; ok {
			c.Release(id, 0, v)
		}
		//打印出错误日志消息
		log.Println(runerrs)
	}
}
