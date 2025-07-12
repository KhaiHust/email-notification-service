package service

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/golibs-starter/golib/log"
)

type IEmailSendingService interface {
	SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) (string, error)
	SendEmailByTask(ctx context.Context, emailRequestID int64) error
}
type EmailSendingService struct {
	emailSendingUsecase    usecase.IEmailSendingUsecase
	getEmailRequestUsecase usecase.IGetEmailRequestUsecase
}

func (e EmailSendingService) SendEmailByTask(ctx context.Context, emailRequestID int64) error {
	emailRequest, err := e.getEmailRequestUsecase.GetEmailRequestByID(ctx, emailRequestID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get email request by id: %d ", emailRequestID), err)
		return err
	}
	return e.emailSendingUsecase.SendEmailByTask(ctx, emailRequest)
}

func (e EmailSendingService) SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) (string, error) {
	reqDto := request.ToEmailSendingRequestDto(req)
	if reqDto == nil {
		log.Error(ctx, "Request is nil")
		return "", nil
	}
	requestID, err := e.emailSendingUsecase.ProcessSendingEmails(ctx, workspaceID, reqDto)
	if err != nil {
		log.Error(ctx, "Error when process sending emails", err)
		return "", err
	}
	return requestID, nil
}

func NewEmailSendingService(
	emailSendingUsecase usecase.IEmailSendingUsecase,
	getEmailRequestUsecase usecase.IGetEmailRequestUsecase,
) IEmailSendingService {
	return &EmailSendingService{
		emailSendingUsecase:    emailSendingUsecase,
		getEmailRequestUsecase: getEmailRequestUsecase,
	}
}
