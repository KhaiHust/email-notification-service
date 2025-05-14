package event

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/event/message"
	"github.com/golibs-starter/golib/web/event"
	"github.com/samber/lo"
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
func NewEventRequestSendingEmail(ctx context.Context, emailRequests []*entity.EmailRequestEntity, req *request.EmailSendingRequestDto) *EventRequestSendingEmail {
	return &EventRequestSendingEmail{
		AbstractEvent: event.NewAbstractEvent(ctx, constant.EmailRequestSendingEvent),
		PayloadData:   EmailRequestEntitiesToMessage(emailRequests, req),
	}
}

func EmailRequestEntitiesToMessage(emailRequests []*entity.EmailRequestEntity, req *request.EmailSendingRequestDto) *message.EmailRequestSendingMessage {
	if len(emailRequests) == 0 || req == nil {
		return nil
	}
	correlationMap := lo.SliceToMap(emailRequests, func(item *entity.EmailRequestEntity) (string, *entity.EmailRequestEntity) {
		return item.CorrelationID, item
	})
	sendDatas := make([]*message.EmailSendData, 0, len(emailRequests))
	for idx, data := range req.Datas {
		correlationKey := fmt.Sprintf("%d_%s", idx, data.To)
		eREntity, ok := correlationMap[correlationKey]
		if !ok {
			continue
		}
		sendDatas = append(sendDatas, &message.EmailSendData{
			EmailRequestID: eREntity.ID, // 🆕 Important: map DB ID
			TrackingID:     eREntity.TrackingID,
			To:             data.To,
			Subject:        data.Subject,
			Body:           data.Body,
		})
	}

	return &message.EmailRequestSendingMessage{
		SendData:      sendDatas,
		TemplateId:    req.TemplateID,
		WorkspaceID:   req.WorkspaceID,
		IntegrationID: req.IntegrationID,
	}
}
