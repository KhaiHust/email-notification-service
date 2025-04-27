package publisher

import (
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/pubsub"
)

type EventPublisherAdapter struct {
}

func (e2 EventPublisherAdapter) Publish(e pubsub.Event) {
	pubsub.Publish(e)
}

func NewEventPublisherAdapter() port.IEventPublisher {
	return &EventPublisherAdapter{}
}
