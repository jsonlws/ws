package beanstalkd

import (
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
