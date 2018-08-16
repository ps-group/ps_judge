package judgeevents

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// EventRouter - routes messages sent to RabbitMQ
type EventRouter struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	lastError error
}

func (router *EventRouter) Error() error {
	return router.lastError
}

// Close - closes router channels and connection
func (router *EventRouter) Close() {
	if router.channel != nil {
		router.channel.Close()
	}
	if router.conn != nil {
		router.conn.Close()
	}
}

// DeclareExchange - declares exchange event source on RabbitMQ
func (router *EventRouter) DeclareExchange(name string) {
	if router.lastError == nil {
		router.lastError = router.channel.ExchangeDeclare(name, "fanout", false, false, false, false, nil)
	}
}

// DeclareQueue - declares message receiver queue
func (router *EventRouter) DeclareQueue(name string) *amqp.Queue {
	if router.lastError == nil {
		var queue amqp.Queue
		queue, router.lastError = router.channel.QueueDeclare(name, false, false, false, false, nil)
		if router.lastError == nil {
			return &queue
		}
	}
	return nil
}

// BindQueue - binds message receiver queue to given exchange event source.
func (router *EventRouter) BindQueue(queue string, exchange string) {
	if router.lastError == nil {
		router.DeclareExchange(exchange)
	}
	if router.lastError == nil {
		router.DeclareQueue(queue)
	}
	if router.lastError == nil {
		router.lastError = router.channel.QueueBind(queue, "", exchange, false, nil)
	}
}

// PublishJSON - publish message with JSON data on given exchange
func (router *EventRouter) PublishJSON(exchange string, value interface{}) {
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

// Consume - creates channel which consumes events from the given queue
func (router *EventRouter) Consume(queue string) <-chan amqp.Delivery {
	var events <-chan amqp.Delivery
	if router.lastError == nil {
		events, router.lastError = router.channel.Consume(queue, "", false, false, false, false, nil)
	}
	return events

}
