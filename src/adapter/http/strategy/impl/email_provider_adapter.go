package strategyAdapterImpl

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/http/strategy"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type EmailProviderAdapter struct {
	providers map[string]strategy.IEmailProviderStrategy
}

func (e EmailProviderAdapter) GetOAuthInfo(ctx context.Context, provider string, code string) (*response.OAuthInfoResponseDto, error) {
	emailProviderStrategy := e.getStrategy(provider)
	if emailProviderStrategy == nil {
		return nil, common.ErrEmailProviderNotFound
	}
	return emailProviderStrategy.GetOAuthInfo(ctx, code)
}

func (e EmailProviderAdapter) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error) {
	emailProviderStrategy := e.getStrategy(provider)
	if emailProviderStrategy == nil {
		return nil, common.ErrEmailProviderNotFound
	}
	return emailProviderStrategy.GetOAuthUrl()
}

func NewEmailProviderAdapter(gmailProviderImpl strategy.IEmailProviderStrategy) port.IEmailProviderPort {
	mapStrategy := make(map[string]strategy.IEmailProviderStrategy)
	mapStrategy[constant.EmailProviderGmail] = gmailProviderImpl
	return &EmailProviderAdapter{
		providers: mapStrategy,
	}
}
func (e EmailProviderAdapter) getStrategy(provider string) strategy.IEmailProviderStrategy {
	return e.providers[provider]
}
