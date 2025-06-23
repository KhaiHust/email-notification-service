package publisher

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib-message-bus/kafka/core"
	"github.com/golibs-starter/golib-message-bus/kafka/relayer"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/pubsub"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type EventPublisherAdapter struct {
	syncProducer   core.SyncProducer
	asyncProducer  core.AsyncProducer
	eventConverter relayer.EventConverter
}

func (e2 EventPublisherAdapter) SyncPublish(ctx context.Context, e pubsub.Event) error {
	message, err := e2.eventConverter.Convert(e)
	if err != nil {
		log.Error(ctx, "Error while converting event to kafka message", err)
		return err
	}
	txn := newrelic.FromContext(ctx)
	log.Info(ctx, fmt.Sprintf("Txn is [%s]", txn.Name()))
	seg := &newrelic.MessageProducerSegment{
		StartTime:       newrelic.StartSegmentNow(txn),
		Library:         "Kafka",
		DestinationType: newrelic.MessageQueue,
		DestinationName: message.Topic,
	}
	partition, offset, err := e2.syncProducer.Send(message)
	defer seg.End()
	if err != nil {
		log.Error(ctx, "Error while sending kafka message", err)
		return err
	}
	log.Info(ctx, "Kafka message sent to partition [%d] with offset [%d]", partition, offset)
	return nil
}

func (e2 EventPublisherAdapter) Publish(e pubsub.Event) {
	pubsub.Publish(e)
}
func NewEventPublisherAdapter(
	syncProducer core.SyncProducer,
	asyncProducer core.AsyncProducer,
	eventConverter relayer.EventConverter,
) port.IEventPublisher {
	return &EventPublisherAdapter{
		syncProducer:   syncProducer,
		asyncProducer:  asyncProducer,
		eventConverter: eventConverter,
	}
}
