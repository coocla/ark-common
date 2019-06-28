package rabbitmq

import (
	"github.com/streadway/amqp"
)

// RabbitQueue 基于RabbitMQ的任务队列
type RabbitQueue struct {
	conn *AMQP
}

// NewRabbitQueue 返回新的rabbitmq连接
func NewRabbitQueue(mqURN string) *RabbitQueue {
	return &RabbitQueue{
		conn: NewQueue(mqURN),
	}
}

// Push 发布消息到指定的Exchange,如果exchange不存在自动创建
func (q *RabbitQueue) Push(exchangeName, exchangeKind string, routingKey string, data []byte) error {
	exchange := &ExchangeContext{
		ExchangeName: exchangeName,
		ExchangeKind: exchangeKind,
		Queues:       []*QueueContext{},
	}
	return q.conn.Publish(exchange, routingKey, string(data))
}

// Listen 监听对应Queue里的消息, Queue不存在则自动创建
func (q *RabbitQueue) Listen(queueName, routingKey string, message chan []byte) error {
	qc := &QueueContext{
		QueueName:  queueName,
		RoutingKey: routingKey,
	}
	return q.conn.Subscribe(qc, message)
}

// BindExchange 将exchange和多个queue绑定在一起
func (q *RabbitQueue) BindExchange(exchangeName, queueName, routingKey string) error {
	exchange := &ExchangeContext{
		ExchangeName: exchangeName,
		ExchangeKind: amqp.ExchangeTopic,
	}
	queue := &QueueContext{
		QueueName:  queueName,
		RoutingKey: routingKey,
	}
	exchange.Queues = []*QueueContext{queue}
	return exchange.BindWithQueues(q.conn)
}
