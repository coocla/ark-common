package rabbitmq

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AMQP struct {
	mu           *sync.RWMutex
	conn         *amqp.Connection
	notifyClose  chan *amqp.Error
	disconnected bool
	Channels     chan *amqp.Channel
	Dial         func() (*amqp.Connection, error)
}

type ExchangeContext struct {
	ExchangeName string
	ExchangeKind string
	Queues       []*QueueContext
}

type QueueContext struct {
	QueueName  string
	RoutingKey string
	AutoAck    bool
}

var MQ *AMQP

// TimeAfter 超时函数
func TimeAfter(d time.Duration) chan int {
	q := make(chan int, 1)
	time.AfterFunc(d, func() {
		q <- 1
	})
	return q
}

func NewQueue(urn string) *AMQP {
	if MQ == nil {
		InitAMQP(urn)
	}
	return MQ
}

// InitAMQP 初始化AMQP的连接
func InitAMQP(urn string) {
	MQ = &AMQP{
		disconnected: true,
		Channels:     make(chan *amqp.Channel, 100),
		notifyClose:  make(chan *amqp.Error),
		Dial: func() (*amqp.Connection, error) {
			return amqp.Dial(urn)
		},
	}
	firstStart := true
	MQ.reConnect(firstStart)
}

func (mq *AMQP) reConnect(firstStart bool) {
	var (
		err  error
		conn *amqp.Connection
	)
	// 重连5次, 每次间隔3秒
	for i := 0; i <= 5; i++ {
		conn, err = mq.Dial()
		if err == nil {
			i = 6
			log.Info("rabbitmq connect success")
			mq.mu.Lock()
			mq.conn = conn
			mq.disconnected = false
			mq.mu.Unlock()
			// 连接关闭监听器
			if firstStart {
				go func() {
					errChan := make(chan *amqp.Error)
					for amqpErr := range conn.NotifyClose(errChan) {
						log.Errorf("rabbitmq disconnected %v, reconnecting", amqpErr)
						mq.mu.Lock()
						mq.disconnected = true
						mq.Channels = make(chan *amqp.Channel, 100)
						mq.mu.Unlock()
						switch {
						case amqpErr.Code == 320:
							mq.reConnect(false)
						case amqpErr.Code == 501:
							mq.reConnect(false)
						case amqpErr.Code == 504:
							mq.reConnect(false)
						}
					}
				}()
			}
		} else {
			if firstStart {
				log.Fatalf("rabbitmq connect failed: %v", err)
			}
		}
		time.Sleep(3 * time.Second)
	}
}

// newChannel 基于底层的conn新建一个channel,多个channel共用一个conn
func (mq *AMQP) newChannel() (channel *amqp.Channel, err error) {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	if mq.disconnected {
		for i := 0; i < 10; i++ {
			if mq.disconnected {
				time.Sleep(3 * time.Second)
				log.Errorf("rabbitmq is disconnected waiting connect")
			} else {
				channel, err = mq.conn.Channel()
				return
			}
		}
		err = errors.New("rabbitmq disconnected")
		return
	}
	channel, err = mq.conn.Channel()
	return
}

// GetChannel 获取一个可用的channel
func (mq *AMQP) GetChannel() (channel *amqp.Channel, err error) {
	for {
		select {
		case channel = <-mq.Channels:
			return
		case <-TimeAfter(time.Second * 1):
			channel, err = mq.newChannel()
			return
		}
	}
}

// ReleaseChannel 释放对应的channel
func (mq *AMQP) ReleaseChannel(channel *amqp.Channel) (closed bool) {
	defer func() {
		err := recover()
		if err != nil {
			mq.Channels = make(chan *amqp.Channel, 100)
			log.Errorf("release channel failed %v", err)
		}
	}()
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	if !mq.disconnected {
		mq.Channels <- channel
	} else {
		channel.Close()
	}
	return
}

// DeclareExchange 创建对应的exchange
func (mq *AMQP) DeclareExchange(e *ExchangeContext) (err error) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	// 名字、类型、是否持久化、是否在没有bind后自动删除、是否作为内部使用、是否要等待队列创建完成
	err = channel.ExchangeDeclare(e.ExchangeName, e.ExchangeKind, true, false, false, false, nil)
	if err != nil {
		log.Errorf("create exchange %s %v", e.ExchangeName, err)
	}
	return
}

// DeleteExchange 删除对应的exchange
func (mq *AMQP) DeleteExchange(e *ExchangeContext) (err error) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	err = channel.ExchangeDelete(e.ExchangeName, false, false)
	return
}

// ExistsExchange 判断对应的exchange是否存在
func (mq *AMQP) ExistsExchange(e *ExchangeContext) (exists bool) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	err = channel.ExchangeDeclarePassive(e.ExchangeName, e.ExchangeKind, true, false, false, false, nil)
	if err == nil {
		exists = true
	}
	return
}

// DeclareQueue 创建对应的Queue
func (mq *AMQP) DeclareQueue(queueName string) (queue amqp.Queue, err error) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	// 名字、是否持久化、是否自动删除、是否本次连接独占使用、是否等待创建结果
	queue, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	return
}

// DeleteQueue 删除对应的Queue
func (mq *AMQP) DeleteQueue(queueName string) (err error) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	_, err = channel.QueueDelete(queueName, false, false, false)
	return
}

// ExistsQueue 判断对应的Queue是否存在
func (mq *AMQP) ExistsQueue(q *QueueContext) (exists bool) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	_, err = channel.QueueDeclarePassive(q.QueueName, true, false, false, false, nil)
	if err == nil {
		exists = true
	}
	return
}

// ExchangeBindWithQueue 将exchange和queue通过对应的routingkey绑定在一起
func (mq *AMQP) ExchangeBindWithQueue(exchangename, routingkey, queuename string) (err error) {
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		return
	}
	err = channel.QueueBind(queuename, routingkey, exchangename, false, nil)
	return
}

// BindWithQueues 自动声明Exchange和Queues并且绑定在一起
func (e *ExchangeContext) BindWithQueues(mq *AMQP) (err error) {
	err = mq.DeclareExchange(e)
	if err != nil {
		return
	}
	for idx := range e.Queues {
		queue := e.Queues[idx]
		_, err = mq.DeclareQueue(queue.QueueName)
		if err != nil {
			return
		}
		err = mq.ExchangeBindWithQueue(e.ExchangeName, queue.RoutingKey, queue.QueueName)
		if err != nil {
			return
		}
	}
	return
}

// Publish 向对应的exchange中发布消息
func (mq *AMQP) Publish(e *ExchangeContext, routingkey string, body string) (err error) {
	err = mq.DeclareExchange(e)
	if err != nil {
		log.Errorf("create exchange %v failed %v", e, err)
		return
	}

	// 如过发生失败，每隔2秒，重试三次
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		log.Errorf("alloc channel failed %v", err)
		return
	}

	loop := 0
	for {
		if err = channel.Publish(e.ExchangeName, routingkey, false, false, amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		}); err != nil {
			if loop > 3 {
				return
			}
			loop++
			log.Errorf("publish message failed %s, retry %d", err, loop)
			time.Sleep(2 * time.Second)
		}
		break
	}
	return
}

// Subscribe 订阅对应queue中的消息
func (mq *AMQP) Subscribe(q *QueueContext, message chan []byte) (err error) {
	// 分配一个channel
	channel, err := mq.GetChannel()
	defer mq.ReleaseChannel(channel)
	if err != nil {
		log.Errorf("alloc channel failed %v", err)
		return
	}

	deliveries, err := channel.Consume(q.QueueName, "", q.AutoAck, false, false, false, nil)
	if err != nil {
		log.Errorf("rabbitmq consume failed, %v", err)
		return
	}
	go func(delivery <-chan amqp.Delivery, message chan []byte) {
		for d := range delivery {
			message <- d.Body
			// 回复本条消息的收到确认
			if !q.AutoAck {
				err := d.Ack(false)
				if err != nil {
					log.Errorf("ack received message failed %v", err)
				}
			}
		}
	}(deliveries, message)
	return
}
