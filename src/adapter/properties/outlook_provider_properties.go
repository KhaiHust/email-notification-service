package properties

import (
	"github.com/golibs-starter/golib/config"
	"golang.org/x/oauth2"
	"strings"
)

type OutlookProviderProperties struct {
	BaseOAuthURL string
	ClientID     string
	RedirectURI  string
	Scopes       string
	ResponseType string
	AccessType   string
	ClientSecret string
	TokenURL     string
	Prompt       string `default:"consent"`
}

func (o OutlookProviderProperties) Prefix() string {
	return "app.services.provider.outlook"
}

func NewOutlookProviderProperties(loader config.Loader) (*OutlookProviderProperties, error) {
	props := &OutlookProviderProperties{}
	err := loader.Bind(props)
	return props, err
}
func NewMicrosoftOAuthConfig(props *OutlookProviderProperties) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     props.ClientID,
		ClientSecret: props.ClientSecret,
		RedirectURL:  props.RedirectURI,
		Scopes:       parseScopes(props.Scopes),
		Endpoint: oauth2.Endpoint{
			AuthURL:  props.BaseOAuthURL,
			TokenURL: props.TokenURL,
		},
	}
}

// Helper to split scope string
func parseScopes(scopes string) []string {
	if scopes == "" {
		return []string{"offline_access", "SMTP.Send", "User.Read"}
	}
	return strings.Fields(scopes)
}
