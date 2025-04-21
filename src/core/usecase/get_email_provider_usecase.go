package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetEmailProviderUseCase interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
	GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error)
}
type GetEmailProviderUseCase struct {
	emailProviderPort           port.IEmailProviderPort
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
}

func (g GetEmailProviderUseCase) GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetEmailProviderByID(ctx, ID)
}

func (g GetEmailProviderUseCase) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error) {
	return g.emailProviderPort.GetOAuthUrl(ctx, provider)
}

func NewGetEmailProviderUseCase(
	emailProviderPort port.IEmailProviderPort,
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
) IGetEmailProviderUseCase {
	return &GetEmailProviderUseCase{
		emailProviderPort:           emailProviderPort,
		emailProviderRepositoryPort: emailProviderRepositoryPort,
	}
}
