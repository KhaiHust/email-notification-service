package handler

import (
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/middleware"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/worker/resource/mapper"
	"github.com/golibs-starter/golib-message-bus/kafka/core"
	"github.com/golibs-starter/golib-message-bus/kafka/relayer"
	"github.com/golibs-starter/golib/log"
)

type EmailSendingRequestHandler struct {
	eventHandlerUsecase usecase.IEventHandlerUsecase
	eventConverter      relayer.EventConverter
}

func (em EmailSendingRequestHandler) HandlerFunc(message *core.ConsumerMessage) {
	var evt event.EventRequestSendingEmail
	if err := em.eventConverter.Restore(message, &evt); err != nil {
		log.Error(evt.Context(), fmt.Sprintf("[EmailSendingRequestHandler] Error unmarshalling message: %v", err))
		return
	}
	ctx := middleware.InitContextWorker(evt.Context())
	if evt.AbstractEvent == nil || evt.AbstractEvent.ApplicationEvent == nil ||
		evt.AbstractEvent.Event != constant.EmailRequestSendingEvent || evt.PayloadData == nil {
		log.Error(ctx, fmt.Sprintf("[EmailSendingRequestHandler] Invalid event: %v", evt))
		return
	}

	payload := evt.PayloadData
	if err := em.eventHandlerUsecase.SendEmailRequestHandler(ctx, payload.IntegrationID, mapper.ToEmailSendingDto(payload)); err != nil {
		log.Error(ctx, fmt.Sprintf("[EmailSendingRequestHandler] Error handling event: %v", err))
		return
	}
	log.Info(ctx, "Successfully processed email sending request")

}

func (em EmailSendingRequestHandler) Close() {
}

func NewEmailSendingRequestHandler(
	eventHandlerUsecase usecase.IEventHandlerUsecase,
	eventConverter relayer.EventConverter,
) core.ConsumerHandler {
	return &EmailSendingRequestHandler{
		eventConverter:      eventConverter,
		eventHandlerUsecase: eventHandlerUsecase,
	}
}
