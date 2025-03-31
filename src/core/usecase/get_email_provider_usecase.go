package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetEmailProviderUseCase interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
}
type GetEmailProviderUseCase struct {
	emailProviderPort port.IEmailProviderPort
}

func (g GetEmailProviderUseCase) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error) {
	return g.emailProviderPort.GetOAuthUrl(ctx, provider)
}

func NewGetEmailProviderUseCase(emailProviderPort port.IEmailProviderPort) IGetEmailProviderUseCase {
	return &GetEmailProviderUseCase{
		emailProviderPort: emailProviderPort,
	}
}
