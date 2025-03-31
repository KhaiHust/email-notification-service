package client

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/strategy"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/golibs-starter/golib/web/client"
	"net/url"
)

type GmailProviderAdapter struct {
	httpClient *client.ContextualHttpClient
	props      *properties.GmailProviderProperties
}

func (g GmailProviderAdapter) GetOAuthUrl() (*response.OAuthUrlResponseDto, error) {
	var params = url.Values{}
	params.Add("client_id", g.props.ClientID)
	params.Add("redirect_uri", g.props.RedirectURI)
	params.Add("response_type", g.props.ResponseType)
	params.Add("scope", g.props.Scope)
	params.Add("access_type", g.props.AccessType)
	oauthUrl := g.props.BaseURL + "?" + params.Encode()
	return &response.OAuthUrlResponseDto{Url: oauthUrl}, nil
}

func NewGmailProviderAdapter(httpClient *client.ContextualHttpClient, props *properties.GmailProviderProperties) strategy.IEmailProviderStrategy {
	return &GmailProviderAdapter{
		httpClient: httpClient,
		props:      props,
	}
}
