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
func (nsc *NatsStreamClient) Connect(clusterID string, clientID string) error {
	nsc.clusterID = clusterID
	nsc.clientID = clientID
	if nsc.PingInterval == 0 {
		nsc.PingInterval = stan.DefaultPingInterval
	}

	if nsc.PingMaxOut == 0 {
		nsc.PingMaxOut = stan.DefaultPingMaxOut
	}

	if nsc.ConnectTimeout == 0 {
		nsc.ConnectTimeout = 2
	}

	return nil
}

// NatsStreamPublish : NatsStream 发布器
type NatsStreamPublish struct {
	NatsStreamClient
	MaxPubNumber int
}

//Connect : 连接Nats服务器
func (nsp *NatsStreamPublish) Connect(clusterID string, clientID string) error {
	nsp.NatsStreamClient.Connect(clusterID, clientID)

	if nsp.MaxPubNumber == 0 {
		nsp.MaxPubNumber = stan.DefaultMaxPubAcksInflight
	}

	err := nsp.autoReConnect()

	if err != nil {
		return err
	}

	return nil
}

// Publish : 发布
func (nsp *NatsStreamPublish) Publish(subject string, data []byte) {
	i := 0
	for {
		err := nsp.c.Publish(subject, data)

		if err == stan.ErrConnectionClosed {
			if nsp.isShutdown {
				logger.Error(nsp.Operator.ID, "nats stream publish closed data not sent subject[%s]:%v", subject, data)
				break
			}
			autoErr := nsp.autoReConnect()
			if autoErr != nil {
				i++
				if i > nsp.AutoReConnectLimt || (nsp.AutoReConnectLimt == 0 && i > defaultAutoReConnectLimt) {
					logger.Error(nsp.Operator.ID,
						"nats stream publish to [%s] data[%v] fail, auto reconnect error:%s",
						subject, data, autoErr.Error())
					break
				}
				logger.Error(nsp.Operator.ID, "nats stream publish auto reconnect [] error :%s", i, autoErr.Error())
				time.Sleep(time.Millisecond * 100)
				continue
			}
		} else if err == nil {
			break
		} else {
			logger.Error(nsp.Operator.ID, "nats stream publish to [%s] data[%v] fail:%s", subject, data, err.Error())
			break
		}

	}
}

// Close : 关闭连接
func (nsp *NatsStreamPublish) Close() {
	if nsp.c != nil {
		nsp.isShutdown = true
		nsp.c.Close()
	}
}

func (nsp *NatsStreamPublish) autoReConnect() error {
	conn, err := stan.Connect(nsp.clusterID,
		nsp.clientID,
		stan.NatsURL(stan.DefaultNatsURL),
		stan.Pings(nsp.PingInterval, nsp.PingMaxOut),
		stan.ConnectWait(time.Duration(nsp.ConnectTimeout)*time.Second),
		stan.MaxPubAcksInflight(nsp.MaxPubNumber),
		stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
			conn.Close()
			logger.Error(nsp.Operator.ID,
				"nats stream conn lost clusterID[%s] clientID[%s] error[%v]",
				nsp.clusterID,
				nsp.clientID,
				err)
		}))

	if err != nil {
		return err
	}
	nsp.c = conn
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
func (nss *NatsStreamSubscribe) Connect(clusterID string, clientID string) error {
	if nss.Quit == nil {
		nss.Quit = make(chan int)
	}

	if nss.MaxSubNumber == 0 {
		nss.MaxSubNumber = stan.DefaultMaxInflight
	}

	if nss.MaxAckSecond == 0 {
		nss.MaxAckSecond = 60
	}

	nss.NatsStreamClient.Connect(clusterID, clientID)
	err := nss.autoReConnect()

	if err != nil {
		return err
	}

	return nil
}

// Subscribe : 以队列订阅 out 通道挂起,需要手动确认Ack包
func (nss *NatsStreamSubscribe) Subscribe(subject string, queueGroup string, durableID string) error {
	awk, _ := time.ParseDuration(fmt.Sprint(nss.MaxAckSecond, "s"))
	sub, err := nss.c.QueueSubscribe(subject, queueGroup, nss.procMessage,
		stan.DurableName(durableID),
		stan.MaxInflight(nss.MaxSubNumber),
		stan.SetManualAckMode(),
		stan.AckWait(awk))

	if err != nil {
		return err
	}
	nss.s = sub
	return nil
}

// Close : 关闭订阅器
func (nss *NatsStreamSubscribe) Close() {
	if nss.s != nil {
		nss.s.Close()
	}

	if nss.c != nil {
		nss.isShutdown = true
		close(nss.Quit)
		nss.c.Close()
		close(nss.Out)
	}
}

func (nss *NatsStreamSubscribe) procMessage(msg *stan.Msg) {
	select {
	case <-nss.Quit:
		return
	default:
	}

	select {
	case <-nss.Quit:
		return
	case nss.Out <- msg:
	}

}

func (nss *NatsStreamSubscribe) autoReConnect() error {
	conn, err := stan.Connect(nss.clusterID,
		nss.clientID,
		stan.NatsURL(stan.DefaultNatsURL),
		stan.Pings(nss.PingInterval, nss.PingMaxOut),
		stan.ConnectWait(time.Duration(nss.ConnectTimeout)*time.Second),
		stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
			conn.Close()
			logger.Error(nss.Operator.ID,
				"nats stream conn lost clusterID[%s] clientID[%s] error[%v]",
				nss.clusterID,
				nss.clientID,
				err)
		}))

	if err != nil {
		return err
	}
	nss.c = conn
	return nil
}
