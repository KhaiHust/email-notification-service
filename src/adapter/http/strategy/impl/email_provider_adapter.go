package strategyAdapterImpl

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/http/strategy"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
	"go.uber.org/fx"
)

type EmailProviderAdapter struct {
	providers map[string]strategy.IEmailProviderStrategy
}

func (e EmailProviderAdapter) GetOAuthByRefreshToken(ctx context.Context, provider *entity.EmailProviderEntity) (*response.OAuthInfoResponseDto, error) {
	emailProviderStrategy := e.getStrategy(provider.Provider)
	if emailProviderStrategy == nil {
		return nil, common.ErrEmailProviderNotFound
	}
	return emailProviderStrategy.GetOAuthByRefreshToken(ctx, provider)
}

func (e EmailProviderAdapter) Send(ctx context.Context, provider *entity.EmailProviderEntity, data *request.EmailDataDto) error {
	emailProviderStrategy := e.getStrategy(provider.Provider)
	if emailProviderStrategy == nil {
		return common.ErrEmailProviderNotFound
	}
	return emailProviderStrategy.SendEmail(ctx, provider, data)
}

func (e EmailProviderAdapter) GetOAuthInfo(ctx context.Context, provider string, code string) (*response.OAuthInfoResponseDto, error) {
	emailProviderStrategy := e.getStrategy(provider)
	if emailProviderStrategy == nil {
		return nil, common.ErrEmailProviderNotFound
	}
	return emailProviderStrategy.GetOAuthInfoByCode(ctx, code)
}

func (e EmailProviderAdapter) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error) {
	emailProviderStrategy := e.getStrategy(provider)
	if emailProviderStrategy == nil {
		return nil, common.ErrEmailProviderNotFound
	}
	return emailProviderStrategy.GetOAuthUrl()
}

type EmailProviderAdapterIn struct {
	fx.In
	Providers []strategy.IEmailProviderStrategy `group:"emailProviderImpl"`
}

func NewEmailProviderAdapter(ins EmailProviderAdapterIn) port.IEmailProviderPort {
	mapStrategy := make(map[string]strategy.IEmailProviderStrategy)
	for _, provider := range ins.Providers {
		mapStrategy[provider.GetType()] = provider
	}
	return &EmailProviderAdapter{
		providers: mapStrategy,
	}
}

func (e EmailProviderAdapter) getStrategy(provider string) strategy.IEmailProviderStrategy {
	return e.providers[provider]
}
