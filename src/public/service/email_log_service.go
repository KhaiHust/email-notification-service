package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
)

type IEmailLogService interface {
	GetLogsByEmailRequestIDAndWorkspaceID(ctx context.Context, emailRequestID int64, workspaceID int64) ([]*entity.EmailLogsEntity, error)
}
type EmailLogService struct {
	getEmailLogUsecase usecase.IGetEmailLogUsecase
}

func (e EmailLogService) GetLogsByEmailRequestIDAndWorkspaceID(ctx context.Context, emailRequestID int64, workspaceID int64) ([]*entity.EmailLogsEntity, error) {
	emailLogs, err := e.getEmailLogUsecase.GetLogsByEmailRequestIDAndWorkspaceID(ctx, emailRequestID, workspaceID)
	if err != nil {
		return nil, err
	}
	return emailLogs, nil
}

func NewEmailLogService(getEmailLogUsecase usecase.IGetEmailLogUsecase) IEmailLogService {
	return &EmailLogService{
		getEmailLogUsecase: getEmailLogUsecase,
	}
}
