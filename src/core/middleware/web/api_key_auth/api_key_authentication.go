package api_key_auth

import (
	"github.com/golibs-starter/golib-security/web/auth/authen"
	"github.com/golibs-starter/golib-security/web/config"
	"github.com/golibs-starter/golib-security/web/filter"
	"go.uber.org/fx"
)

func ApiKeyAuthFilterOpt() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  "authentication_filter",
		Target: NewApiKeyAuthFilter,
	})
}

type ApiKeyAuthFilterIn struct {
	fx.In
	SecurityProperties  *config.HttpSecurityProperties
	AuthProviderManager *authen.ProviderManager
}

func NewApiKeyAuthFilter(in ApiKeyAuthFilterIn) (filter.AuthenticationFilter, error) {

	apiKeyFilter, err := ApiKeyAuthFilter()
	if err != nil {
		return nil, err
	}
	return apiKeyFilter, nil
}
