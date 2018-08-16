package judgeevents

import (
	"github.com/streadway/amqp"
)

// NewEventRouter - creates message router
func NewEventRouter(socket string) *EventRouter {
	var router EventRouter
	router.conn, router.lastError = amqp.Dial(socket)
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
