package library

import (
	"fmt"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"

	stan "github.com/nats-io/go-nats-streaming"
)

const (
	defaultAutoReConnectLimt = 3
)

// NatsStreamClient :
type NatsStreamClient struct {
	c          stan.Conn
	isShutdown bool
	clusterID  string
	clientID   string

	AutoReConnectLimt int //自动重连最大次数
	Operator          *actor.PID
	PingInterval      int
	PingMaxOut        int
	ConnectTimeout    int
}

// Connect : xx
func (slf *NatsStreamClient) Connect(clusterID string, clientID string) error {
	slf.clusterID = clusterID
	slf.clientID = clientID
	if slf.PingInterval == 0 {
		slf.PingInterval = stan.DefaultPingInterval
	}

	if slf.PingMaxOut == 0 {
		slf.PingMaxOut = stan.DefaultPingMaxOut
	}

	if slf.ConnectTimeout == 0 {
		slf.ConnectTimeout = 2
	}

	return nil
}

// NatsStreamPublish : NatsStream 发布器
type NatsStreamPublish struct {
	NatsStreamClient
	MaxPubNumber int
}

//Connect : 连接Nats服务器
func (slf *NatsStreamPublish) Connect(clusterID string, clientID string) error {
	slf.NatsStreamClient.Connect(clusterID, clientID)

	if slf.MaxPubNumber == 0 {
		slf.MaxPubNumber = stan.DefaultMaxPubAcksInflight
	}

	err := slf.autoReConnect()

	if err != nil {
		return err
	}

	return nil
}

// Publish : 发布
func (slf *NatsStreamPublish) Publish(subject string, data []byte) {
	i := 0
	for {
		err := slf.c.Publish(subject, data)

		if err == stan.ErrConnectionClosed {
			if slf.isShutdown {
				logger.Error(slf.Operator.ID, "nats stream publish closed data not sent subject[%s]:%v", subject, data)
				break
			}
			autoErr := slf.autoReConnect()
			if autoErr != nil {
				i++
				if i > slf.AutoReConnectLimt || (slf.AutoReConnectLimt == 0 && i > defaultAutoReConnectLimt) {
					logger.Error(slf.Operator.ID,
						"nats stream publish to [%s] data[%v] fail, auto reconnect error:%s",
						subject, data, autoErr.Error())
					break
				}
				logger.Error(slf.Operator.ID, "nats stream publish auto reconnect [] error :%s", i, autoErr.Error())
				time.Sleep(time.Millisecond * 100)
				continue
			}
		} else if err == nil {
			break
		} else {
			logger.Error(slf.Operator.ID, "nats stream publish to [%s] data[%v] fail:%s", subject, data, err.Error())
			break
		}

	}
}

// Close : 关闭连接
func (slf *NatsStreamPublish) Close() {
	if slf.c != nil {
		slf.isShutdown = true
		slf.c.Close()
	}
}

func (slf *NatsStreamPublish) autoReConnect() error {
	conn, err := stan.Connect(slf.clusterID,
		slf.clientID,
		stan.NatsURL(stan.DefaultNatsURL),
		stan.Pings(slf.PingInterval, slf.PingMaxOut),
		stan.ConnectWait(time.Duration(slf.ConnectTimeout)*time.Second),
		stan.MaxPubAcksInflight(slf.MaxPubNumber),
		stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
			conn.Close()
			logger.Error(slf.Operator.ID,
				"nats stream conn lost clusterID[%s] clientID[%s] error[%v]",
				slf.clusterID,
				slf.clientID,
				err)
		}))

	if err != nil {
		return err
	}
	slf.c = conn
	return nil
}

// NatsStreamSubscribe : 订阅
type NatsStreamSubscribe struct {
	s stan.Subscription
	NatsStreamClient
	MaxAckSecond int
	MaxSubNumber int
	Out          chan *stan.Msg
	Quit         chan int
}

// Connect : xxx
func (slf *NatsStreamSubscribe) Connect(clusterID string, clientID string) error {
	if slf.Quit == nil {
		slf.Quit = make(chan int)
	}

	if slf.MaxSubNumber == 0 {
		slf.MaxSubNumber = stan.DefaultMaxInflight
	}

	if slf.MaxAckSecond == 0 {
		slf.MaxAckSecond = 60
	}

	slf.NatsStreamClient.Connect(clusterID, clientID)
	err := slf.autoReConnect()

	if err != nil {
		return err
	}

	return nil
}

// Subscribe : 以队列订阅 out 通道挂起,需要手动确认Ack包
func (slf *NatsStreamSubscribe) Subscribe(subject string, queueGroup string, durableID string) error {
	awk, _ := time.ParseDuration(fmt.Sprint(slf.MaxAckSecond, "s"))
	sub, err := slf.c.QueueSubscribe(subject, queueGroup, slf.procMessage,
		stan.DurableName(durableID),
		stan.MaxInflight(slf.MaxSubNumber),
		stan.SetManualAckMode(),
		stan.AckWait(awk))

	if err != nil {
		return err
	}
	slf.s = sub
	return nil
}

// Close : 关闭订阅器
func (slf *NatsStreamSubscribe) Close() {
	if slf.s != nil {
		slf.s.Close()
	}

	if slf.c != nil {
		slf.isShutdown = true
		close(slf.Quit)
		slf.c.Close()
		close(slf.Out)
	}
}

func (slf *NatsStreamSubscribe) procMessage(msg *stan.Msg) {
	select {
	case <-slf.Quit:
		return
	default:
	}

	select {
	case <-slf.Quit:
		return
	case slf.Out <- msg:
	}

}

func (slf *NatsStreamSubscribe) autoReConnect() error {
	conn, err := stan.Connect(slf.clusterID,
		slf.clientID,
		stan.NatsURL(stan.DefaultNatsURL),
		stan.Pings(slf.PingInterval, slf.PingMaxOut),
		stan.ConnectWait(time.Duration(slf.ConnectTimeout)*time.Second),
		stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
			conn.Close()
			logger.Error(slf.Operator.ID,
				"nats stream conn lost clusterID[%s] clientID[%s] error[%v]",
				slf.clusterID,
				slf.clientID,
				err)
		}))

	if err != nil {
		return err
	}
	slf.c = conn
	return nil
}
