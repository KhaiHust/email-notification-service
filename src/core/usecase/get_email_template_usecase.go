package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetEmailTemplateUseCase interface {
	GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error)
}
type GetEmailTemplateUseCase struct {
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort
}

func (e GetEmailTemplateUseCase) GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error) {
	return e.emailTemplateRepositoryPort.GetTemplateByID(ctx, ID)
}

func NewEmailTemplateUseCase(
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort,
) IGetEmailTemplateUseCase {
	return &GetEmailTemplateUseCase{
		emailTemplateRepositoryPort: emailTemplateRepositoryPort,
	}
}
