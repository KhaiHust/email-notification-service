package strategy

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
)

type IEmailProviderStrategy interface {
	GetType() string
	GetOAuthUrl() (*response.OAuthUrlResponseDto, error)
	GetOAuthInfoByCode(ctx context.Context, code string) (*response.OAuthInfoResponseDto, error)
	SendEmail(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity, emailData *request.EmailDataDto) error
	GetOAuthByRefreshToken(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*response.OAuthInfoResponseDto, error)
}
