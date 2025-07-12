package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IGetEmailLogUsecase interface {
	GetLogsByEmailRequestIDAndWorkspaceID(ctx context.Context, emailRequestID int64, workspaceID int64) ([]*entity.EmailLogsEntity, error)
}
type GetEmailLogUsecase struct {
	emailLogRepositoryPort port.IEmailLogRepositoryPort
}

func (g GetEmailLogUsecase) GetLogsByEmailRequestIDAndWorkspaceID(ctx context.Context, emailRequestID int64, workspaceID int64) ([]*entity.EmailLogsEntity, error) {
	emailLogs, err := g.emailLogRepositoryPort.GetLogsByEmailRequestIDAndWorkspaceID(ctx, emailRequestID, workspaceID)
	if err != nil {
		log.Error(ctx, "Error when get email logs by email request id", err)
		return nil, err
	}
	return emailLogs, nil
}

func NewGetEmailLogUsecase(emailLogRepositoryPort port.IEmailLogRepositoryPort) IGetEmailLogUsecase {
	return &GetEmailLogUsecase{
		emailLogRepositoryPort: emailLogRepositoryPort,
	}
}
