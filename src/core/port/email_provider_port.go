package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
)

type IEmailProviderPort interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
}
