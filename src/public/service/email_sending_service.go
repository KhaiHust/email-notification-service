package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/golibs-starter/golib/log"
)

type IEmailSendingService interface {
	SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) error
}
type EmailSendingService struct {
	emailSendingUsecase usecase.IEmailSendingUsecase
}

func (e EmailSendingService) SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) error {
	reqDto := request.ToEmailSendingRequestDto(req)
	if reqDto == nil {
		log.Error(ctx, "Request is nil")
		return nil
	}
	if err := e.emailSendingUsecase.ProcessSendingEmails(ctx, workspaceID, reqDto); err != nil {
		log.Error(ctx, "Error when process sending emails", err)
	}
	return nil
}

func NewEmailSendingService(
	emailSendingUsecase usecase.IEmailSendingUsecase,
) IEmailSendingService {
	return &EmailSendingService{
		emailSendingUsecase: emailSendingUsecase,
	}
}
