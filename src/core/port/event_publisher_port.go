package port

import (
	"context"
	"github.com/golibs-starter/golib/pubsub"
)

type IEventPublisher interface {
	Publish(e pubsub.Event)
	SyncPublish(ctx context.Context, e pubsub.Event) error
}
