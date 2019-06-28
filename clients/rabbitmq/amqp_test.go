package rabbitmq

import (
	"testing"

	"github.com/streadway/amqp"
)

func init() {
	InitCfg("../cfg.json")
	InitAMQP()
}
func TestDeclareDeleteExchange(t *testing.T) {
	exchange := &ExchangeContext{
		ExchangeName: "test2",
		ExchangeKind: amqp.ExchangeFanout,
	}
	err := MQ.DeclareExchange(exchange)
	if err != nil {
		t.Errorf("Declare Exchange failed %v", err)
	}
	err = MQ.DeleteExchange(exchange)
	if err != nil {
		t.Errorf("Delete Exchange failed %v", err)
	}
}

func TestDeclareDeleteQueue(t *testing.T) {
	queue := &QueueContext{
		QueueName:  "test",
		RoutingKey: "test",
	}
	_, err := MQ.DeclareQueue(queue)
	if err != nil {
		t.Errorf("Declare Queue failed %v", err)
	}
	err = MQ.DeleteQueue(queue)
	if err != nil {
		t.Errorf("Delete Queue failed %v", err)
	}
}

func TestPublishMessage(t *testing.T) {
	exchange := &ExchangeContext{
		ExchangeName: "test2",
		ExchangeKind: amqp.ExchangeFanout,
		Queues: []*QueueContext{
			&QueueContext{
				QueueName:  "test1",
				RoutingKey: "test",
			},
			&QueueContext{
				QueueName:  "test2",
				RoutingKey: "test",
			},
		},
	}
	err := MQ.Publish(exchange, "test", "test")
	if err != nil {
		t.Errorf("Publish Message failed %v", err.Error())
	}
}

func TestSubscribeMessage(t *testing.T) {
	queue := &QueueContext{
		QueueName:  "test2",
		RoutingKey: "test",
	}
	reply := make(chan []byte)
	err := MQ.Subscribe(queue, reply)
	if err != nil {
		t.Errorf("Subscribe Message failed %v", err)
	}
	exchange := &ExchangeContext{
		ExchangeName: "test2",
		ExchangeKind: amqp.ExchangeFanout,
		Queues:       []*QueueContext{queue},
	}
	message := "test"
	err = MQ.Publish(exchange, queue.RoutingKey, message)
	if err != nil {
		t.Errorf("Publish Message failed %v", err)
	}

	t.Run("TestReceiveMessage", func(t *testing.T) {
		msg := <-reply
		if string(msg) != message {
			t.Errorf("Receive Message failed %v", err)
		}
	})
}
