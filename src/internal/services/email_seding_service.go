package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/internal/resources/request"
	"github.com/KhaiHust/email-notification-service/internal/resources/response"
	"github.com/golibs-starter/golib/log"
)

type IEmailSendingService interface {
	SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) (*response.EmailSendingResponse, error)
	SendEmailByTask(ctx context.Context, emailRequestID int64) error
}
type EmailSendingService struct {
	emailSendingUsecase     usecase.IEmailSendingUsecase
	getEmailProviderUseCase usecase.IGetEmailProviderUseCase
	getEmailRequestUsecase  usecase.IGetEmailRequestUsecase
}

func (e EmailSendingService) SendEmailByTask(ctx context.Context, emailRequestID int64) error {
	emailRequest, err := e.getEmailRequestUsecase.GetEmailRequestByID(ctx, emailRequestID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get email request by id: %d ", emailRequestID), err)
		return err
	}
	return e.emailSendingUsecase.SendEmailByTask(ctx, emailRequest)
}

func (e EmailSendingService) SendEmailRequest(ctx context.Context, workspaceID int64, req *request.EmailSendingRequest) (*response.EmailSendingResponse, error) {
	provider, err := e.getEmailProviderUseCase.GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx, req.Provider, workspaceID, req.Environment)
	if err != nil {
		log.Error(ctx, "Failed to get email provider", err)
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, common.ErrProviderNotFoundOrForbidden
		}
		return nil, err
	}
	reqDto := request.ToEmailSendingRequestDto(req)
	reqDto.Provider = provider
	requestID, err := e.emailSendingUsecase.ProcessSendingEmails(ctx, workspaceID, reqDto)
	if err != nil {
		return nil, err
	}
	return &response.EmailSendingResponse{RequestID: requestID}, nil
}

func NewEmailSendingService(
	emailSendingUsecase usecase.IEmailSendingUsecase,
	getEmailProviderUseCase usecase.IGetEmailProviderUseCase,
	getEmailRequestUsecase usecase.IGetEmailRequestUsecase,
) IEmailSendingService {
	return &EmailSendingService{
		emailSendingUsecase:     emailSendingUsecase,
		getEmailProviderUseCase: getEmailProviderUseCase,
		getEmailRequestUsecase:  getEmailRequestUsecase,
	}
}
