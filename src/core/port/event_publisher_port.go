package port

import "github.com/golibs-starter/golib/pubsub"

type IEventPublisher interface {
	Publish(e pubsub.Event)
}
