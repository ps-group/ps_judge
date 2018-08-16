package judgeevents

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

const (
	// ExchangeBuildFinished - exchange name for BuildFinished event
	ExchangeBuildFinished = "psjudge-build-finished"
)

// BuildFinishedCallback - callback which receives BuildFinished event
type BuildFinishedCallback func(event BuildFinishedEvent)

// BuildFinishedEvent - represents info attached to BuildFinished event
type BuildFinishedEvent struct {
	Key     string `json:"key"`
	Succeed bool   `json:"succeed"`
}

// BuilderEvents - allows to publish and subscribe to the Builder events.
type BuilderEvents interface {
	Error() error
	PublishBuildFinished(event BuildFinishedEvent)
	ConsumeBuildFinished(queue string, cb BuildFinishedCallback)
	Close()
}

type builderEventsImpl struct {
	router  *EventRouter
	closing chan struct{}
	started bool
}

// NewBuilderEvents - creates new BuilderEvents router
func NewBuilderEvents(socket string) BuilderEvents {
	events := new(builderEventsImpl)
	events.router = NewEventRouter(socket)
	return events
}

func (impl *builderEventsImpl) Error() error {
	return impl.router.Error()
}

func (impl *builderEventsImpl) PublishBuildFinished(event BuildFinishedEvent) {
	impl.router.DeclareExchange(ExchangeBuildFinished)
	impl.router.PublishJSON(ExchangeBuildFinished, event)
}

func (impl *builderEventsImpl) ConsumeBuildFinished(queue string, cb BuildFinishedCallback) {
	impl.router.BindQueue(queue, ExchangeBuildFinished)

	channel := impl.router.Consume(queue)
	if channel != nil {
		impl.closing = make(chan struct{})
		impl.started = true

		go func() {
			active := true
			for active {
				var event BuildFinishedEvent
				select {
				case delivery := <-channel:
					err := json.Unmarshal(delivery.Body, &event)
					if err != nil {
						logrus.WithField("error", err).Error("cannot decode router message")
					}
					cb(event)
				case <-impl.closing:
					impl.closing <- struct{}{}
					active = false
				}
			}
		}()
	}
}

func (impl *builderEventsImpl) Close() {
	if impl.started {
		impl.closing <- struct{}{}
		<-impl.closing
		impl.started = false
	}
}
