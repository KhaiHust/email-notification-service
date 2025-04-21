package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
)

type IEventHandlerUsecase interface {
	SendEmailRequestHandler(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error
}
type EventHandlerUsecase struct {
	emailSendingUsecase IEmailSendingUsecase
}

func (e EventHandlerUsecase) SendEmailRequestHandler(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error {
	return e.emailSendingUsecase.SendBatches(ctx, providerID, req)
}

func NewEventHandlerUsecase(
	emailSendingUsecase IEmailSendingUsecase,
) IEventHandlerUsecase {
	return &EventHandlerUsecase{
		emailSendingUsecase: emailSendingUsecase,
	}
}
