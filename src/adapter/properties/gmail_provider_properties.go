package properties

import (
	"github.com/golibs-starter/golib/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"strings"
)

type GmailProviderProperties struct {
	BaseOAuthURL string
	ClientID     string
	RedirectURI  string
	Scopes       string
	ResponseType string
	AccessType   string
	ClientSecret string
}

func (g GmailProviderProperties) Prefix() string {
	return "app.services.provider.gmail"
}

func NewGmailProviderProperties(loader config.Loader) (*GmailProviderProperties, error) {
	prop := &GmailProviderProperties{}
	err := loader.Bind(prop)
	return prop, err
}

type GoogleOAuthConfig struct {
	props *GmailProviderProperties
}

func NewGoogleOAuthConfig(props *GmailProviderProperties) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     props.ClientID,
		RedirectURL:  props.RedirectURI,
		ClientSecret: props.ClientSecret,
		Scopes:       strings.Split(props.Scopes, " "),
		Endpoint:     google.Endpoint,
	}
}
