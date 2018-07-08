package main

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

const (
	RabbitMQSocket        = "amqp://guest:guest@localhost:5672/"
	ExchangeBuildFinished = "psjudge-build-finished"
)

type BuildFinishedEvent struct {
	Key     string `json:"key"`
	Succeed bool   `json:"succeed"`
}

type MessageRouter struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	lastError error
}

func (router *MessageRouter) Error() error {
	return router.lastError
}

func (router *MessageRouter) Close() {
	if router.channel != nil {
		router.channel.Close()
	}
	if router.conn != nil {
		router.conn.Close()
	}
}

func NewMessageRouter() MessageRouter {
	var router MessageRouter
	router.conn, router.lastError = amqp.Dial(RabbitMQSocket)
	if router.lastError == nil {
		defer func() {
			if router.channel == nil {
				router.conn.Close()
				router.conn = nil
			}
		}()
		router.channel, router.lastError = router.conn.Channel()
	}
	return router
}

func (router *MessageRouter) DeclareExchange(name string) {
	if router.lastError == nil {
		router.lastError = router.channel.ExchangeDeclare(name, "fanout", false, false, false, false, nil)
	}
}

func (router *MessageRouter) DeclareQueue(name string) *amqp.Queue {
	if router.lastError == nil {
		var queue amqp.Queue
		queue, router.lastError = router.channel.QueueDeclare(name, false, false, false, false, nil)
		if router.lastError != nil {
			return &queue
		}
	}
	return nil
}

func (router *MessageRouter) BindQueue(queue string, exchange string) {
	if router.lastError != nil {
		router.DeclareExchange(exchange)
	}
	if router.lastError != nil {
		router.DeclareQueue(queue)
	}
	if router.lastError != nil {
		router.lastError = router.channel.QueueBind(queue, "", exchange, false, nil)
	}
}

func (router *MessageRouter) PublishJSON(exchange string, value interface{}) {
	if router.lastError == nil {
		var body []byte
		body, router.lastError = json.Marshal(value)
		if router.lastError == nil {
			router.lastError = router.channel.Publish(exchange, "", false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		}
	}
}

func (router *MessageRouter) Consume(queue string) <-chan amqp.Delivery {
	var messages <-chan amqp.Delivery
	if router.lastError == nil {
		messages, router.lastError = router.channel.Consume(queue, "", false, false, false, false, nil)
	}
	return messages
}
