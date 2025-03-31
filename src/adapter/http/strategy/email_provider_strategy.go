package strategy

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
)

type IEmailProviderStrategy interface {
	GetOAuthUrl() (*response.OAuthUrlResponseDto, error)
	GetOAuthInfo(ctx context.Context, code string) (*response.OAuthInfoResponseDto, error)
}
