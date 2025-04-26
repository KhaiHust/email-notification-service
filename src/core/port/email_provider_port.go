package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
)

type IEmailProviderPort interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
	GetOAuthInfo(ctx context.Context, provider string, code string) (*response.OAuthInfoResponseDto, error)
	Send(ctx context.Context, provider *entity.EmailProviderEntity, data *request.EmailDataDto) error
	GetOAuthByRefreshToken(ctx context.Context, provider *entity.EmailProviderEntity) (*response.OAuthInfoResponseDto, error)
}
