package main

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

const (
	// ExchangeBuildFinished - exchange name for BuildFinished event
	ExchangeBuildFinished = "psjudge-build-finished"
)

// BuildFinishedEvent - represents info attached to BuildFinished event
type BuildFinishedEvent struct {
	Key     string `json:"key"`
	Succeed bool   `json:"succeed"`
}

// MessageRouter - routes messages sent to RabbitMQ
type MessageRouter struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	lastError error
}

func (router *MessageRouter) Error() error {
	return router.lastError
}

// Close - closes router channels and connection
func (router *MessageRouter) Close() {
	if router.channel != nil {
		router.channel.Close()
	}
	if router.conn != nil {
		router.conn.Close()
	}
}

// MessageRouterFactory - creates new message routers
type MessageRouterFactory interface {
	NewMessageRouter() *MessageRouter
}

type messageRouterFactoryImpl struct {
	socket string
}

// NewMessageRouterFactory - creates factory that creates routers on given socket
func NewMessageRouterFactory(socket string) MessageRouterFactory {
	f := new(messageRouterFactoryImpl)
	f.socket = socket
	return f
}

// NewMessageRouter - creates message router
func (f *messageRouterFactoryImpl) NewMessageRouter() *MessageRouter {
	var router MessageRouter
	router.conn, router.lastError = amqp.Dial(f.socket)
	if router.lastError == nil {
		defer func() {
			if router.channel == nil {
				router.conn.Close()
				router.conn = nil
			}
		}()
		router.channel, router.lastError = router.conn.Channel()
	}
	return &router
}

// DeclareExchange - declares exchange on RabbitMQ
func (router *MessageRouter) DeclareExchange(name string) {
	if router.lastError == nil {
		router.lastError = router.channel.ExchangeDeclare(name, "fanout", false, false, false, false, nil)
	}
}

func (router *MessageRouter) declareQueue(name string) *amqp.Queue {
	if router.lastError == nil {
		var queue amqp.Queue
		queue, router.lastError = router.channel.QueueDeclare(name, false, false, false, false, nil)
		if router.lastError != nil {
			return &queue
		}
	}
	return nil
}

func (router *MessageRouter) bindQueue(queue string, exchange string) {
	if router.lastError != nil {
		router.DeclareExchange(exchange)
	}
	if router.lastError != nil {
		router.declareQueue(queue)
	}
	if router.lastError != nil {
		router.lastError = router.channel.QueueBind(queue, "", exchange, false, nil)
	}
}

// PublishJSON - publish message with JSON data on given exchange
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

func (router *MessageRouter) consume(queue string) <-chan amqp.Delivery {
	var messages <-chan amqp.Delivery
	if router.lastError == nil {
		messages, router.lastError = router.channel.Consume(queue, "", false, false, false, false, nil)
	}
	return messages
}
