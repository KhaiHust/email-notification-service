package event

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/event/message"
	"github.com/golibs-starter/golib/web/event"
)

type EventRequestSendingEmail struct {
	*event.AbstractEvent
	PayloadData *message.EmailRequestSendingMessage `json:"payload"`
}

func (a EventRequestSendingEmail) Payload() interface{} {
	return a.PayloadData
}
func (a EventRequestSendingEmail) String() string {
	return a.ToString(a)
}
func NewEventRequestSendingEmail(ctx context.Context, emailRequests []*entity.EmailRequestEntity) *EventRequestSendingEmail {
	return &EventRequestSendingEmail{
		AbstractEvent: event.NewAbstractEvent(ctx, constant.EmailRequestSendingEvent),
		PayloadData:   EmailRequestEntitiesToMessage(emailRequests),
	}
}

func EmailRequestEntitiesToMessage(emailRequests []*entity.EmailRequestEntity) *message.EmailRequestSendingMessage {
	sendDatas := make([]*message.EmailSendData, 0, len(emailRequests))
	for _, eREntity := range emailRequests {
		sendDatas = append(sendDatas, &message.EmailSendData{
			EmailRequestID: eREntity.ID,
			TrackingID:     eREntity.TrackingID,
			SendAt:         eREntity.SendAt,
			IsRetry:        eREntity.IsRetry,
		})
	}

	return &message.EmailRequestSendingMessage{
		SendData: sendDatas,
	}
}
