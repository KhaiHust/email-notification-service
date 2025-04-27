package event

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/event/message"
	"github.com/golibs-starter/golib/web/event"
)

type EventEmailRequestSync struct {
	*event.AbstractEvent
	PayloadData *message.EmailRequestSyncMessage `json:"payload"`
}

func (a EventEmailRequestSync) Payload() interface{} {
	return a.PayloadData
}
func (a EventEmailRequestSync) String() string {
	return a.ToString(a)
}
func NewEventEmailRequestSync(ctx context.Context, emailRequest *entity.EmailRequestEntity) *EventEmailRequestSync {
	return &EventEmailRequestSync{
		AbstractEvent: event.NewAbstractEvent(ctx, constant.EmailRequestSyncEvent),
		PayloadData:   EmailRequestEntityToMessage(emailRequest),
	}
}

func EmailRequestEntityToMessage(request *entity.EmailRequestEntity) *message.EmailRequestSyncMessage {
	if request == nil {
		return nil
	}
	return &message.EmailRequestSyncMessage{
		EmailRequestID: request.ID,
		Status:         request.Status,
		ErrorMessage:   request.ErrorMessage,
		SentAt:         request.SentAt,
	}
}
func MessageToEmailRequestEntity(msg *message.EmailRequestSyncMessage) *entity.EmailRequestEntity {
	if msg == nil {
		return nil
	}
	return &entity.EmailRequestEntity{
		BaseEntity: entity.BaseEntity{
			ID: msg.EmailRequestID,
		},
		Status:       msg.Status,
		ErrorMessage: msg.ErrorMessage,
		SentAt:       msg.SentAt,
	}
}
