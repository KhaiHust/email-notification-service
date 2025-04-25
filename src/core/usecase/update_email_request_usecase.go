package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IUpdateEmailRequestUsecase interface {
	UpdateStatusByBatches(ctx context.Context, emailRequestEntity []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
}
type UpdateEmailRequestUsecase struct {
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
}

func (u UpdateEmailRequestUsecase) UpdateStatusByBatches(ctx context.Context, emailRequestEntity []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	//TODO implement me
	panic("implement me")
}

func NewUpdateEmailRequestUsecase(emailRequestRepositoryPort port.IEmailRequestRepositoryPort) IUpdateEmailRequestUsecase {
	return &UpdateEmailRequestUsecase{
		emailRequestRepositoryPort: emailRequestRepositoryPort,
	}
}
