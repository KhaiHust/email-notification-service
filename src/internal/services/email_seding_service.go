package services

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/internal/resources/request"
)

type IEmailSendingService interface {
	SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) (interface{}, error)
}
type EmailSendingService struct {
	emailSendingUsecase usecase.IEmailSendingUsecase
}

func (e EmailSendingService) SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func NewEmailSendingService(
	emailSendingUsecase usecase.IEmailSendingUsecase,
) IEmailSendingService {
	return &EmailSendingService{
		emailSendingUsecase: emailSendingUsecase,
	}
}
