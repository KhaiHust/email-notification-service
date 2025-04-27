package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/golibs-starter/golib/log"
)

type IEventHandlerUsecase interface {
	SendEmailRequestHandler(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error
	SyncEmailRequestHandler(ctx context.Context, emailRequest *entity.EmailRequestEntity) error
}
type EventHandlerUsecase struct {
	emailSendingUsecase       IEmailSendingUsecase
	updateEmailRequestUsecase IUpdateEmailRequestUsecase
	getEmailRequestUsecase    IGetEmailRequestUsecase
}

func (e EventHandlerUsecase) SyncEmailRequestHandler(ctx context.Context, emailRequest *entity.EmailRequestEntity) error {
	emailRequestEntity, err := e.getEmailRequestUsecase.GetEmailRequestByID(ctx, emailRequest.ID)
	if err != nil {
		log.Error(ctx, "Error when get email request by id", err)
		return err
	}
	emailRequestEntity.Status = emailRequest.Status
	emailRequestEntity.ErrorMessage = emailRequest.ErrorMessage
	emailRequestEntity.SentAt = emailRequest.SentAt
	if _, err := e.updateEmailRequestUsecase.UpdateEmailRequestByID(ctx, emailRequestEntity); err != nil {
		log.Error(ctx, "Error when update email request by id", err)
		return err
	}
	return nil
}

func (e EventHandlerUsecase) SendEmailRequestHandler(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error {
	return e.emailSendingUsecase.SendBatches(ctx, providerID, req)
}

func NewEventHandlerUsecase(
	emailSendingUsecase IEmailSendingUsecase,
	updateEmailRequestUsecase IUpdateEmailRequestUsecase,
	getEmailRequestUsecase IGetEmailRequestUsecase,
) IEventHandlerUsecase {
	return &EventHandlerUsecase{
		emailSendingUsecase:       emailSendingUsecase,
		updateEmailRequestUsecase: updateEmailRequestUsecase,
		getEmailRequestUsecase:    getEmailRequestUsecase,
	}
}
