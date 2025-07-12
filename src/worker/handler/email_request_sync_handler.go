package handler

import (
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/middleware"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/golibs-starter/golib-message-bus/kafka/core"
	"github.com/golibs-starter/golib-message-bus/kafka/relayer"
	"github.com/golibs-starter/golib/log"
)

type EmailRequestSyncHandler struct {
	eventHandlerUsecase usecase.IEventHandlerUsecase
	eventConverter      relayer.EventConverter
}

func (em EmailRequestSyncHandler) Close() {

}

func (em EmailRequestSyncHandler) HandlerFunc(message *core.ConsumerMessage) {
	var evt event.EventEmailRequestSync
	if err := em.eventConverter.Restore(message, &evt); err != nil {
		log.Error(fmt.Sprintf("[EmailRequestSyncHandler] Error unmarshalling message: %v, with "+
			"event [%v]", err, constant.EmailRequestSyncEvent))
		return
	}

	ctx := middleware.InitContextWorker(evt.Context())

	if evt.AbstractEvent == nil || evt.AbstractEvent.ApplicationEvent == nil ||
		evt.AbstractEvent.Event != constant.EmailRequestSyncEvent || evt.PayloadData == nil {
		log.Error(ctx, fmt.Sprintf("[EmailRequestSyncHandler] Invalid event: %v", evt))
		return
	}
	emailRequest := event.MessageToEmailRequestEntity(evt.PayloadData)
	if emailRequest == nil {
		log.Error(ctx, "[EmailRequestSyncHandler] Invalid payload data")
		return
	}
	if err := em.eventHandlerUsecase.SyncEmailRequestHandler(ctx, emailRequest); err != nil {
		log.Error(ctx, "[EmailRequestSyncHandler] Error when syncing email request", err)
		return
	}
	log.Info(ctx, "[EmailRequestSyncHandler] Sync email request successfully with id: %d", emailRequest.ID)
}
func NewEmailRequestSyncHandler(
	eventHandlerUsecase usecase.IEventHandlerUsecase,
	eventConverter relayer.EventConverter,
) core.ConsumerHandler {
	return &EmailRequestSyncHandler{
		eventHandlerUsecase: eventHandlerUsecase,
		eventConverter:      eventConverter,
	}
}
