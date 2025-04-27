package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetEmailRequestUsecase interface {
	GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error)
}
type GetEmailRequestUsecase struct {
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
}

func (g GetEmailRequestUsecase) GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error) {
	return g.emailRequestRepositoryPort.GetEmailRequestByID(ctx, emailRequestID)
}

func NewGetEmailRequestUsecase(
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
) IGetEmailRequestUsecase {
	return &GetEmailRequestUsecase{
		emailRequestRepositoryPort: emailRequestRepositoryPort,
	}
}
