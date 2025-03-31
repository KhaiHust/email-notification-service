package strategy

import "github.com/KhaiHust/email-notification-service/core/entity/dto/response"

type IEmailProviderStrategy interface {
	GetOAuthUrl() (*response.OAuthUrlResponseDto, error)
}
