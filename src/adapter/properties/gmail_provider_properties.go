package properties

import (
	"github.com/golibs-starter/golib/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GmailProviderProperties struct {
	BaseOAuthURL string
	ClientID     string
	RedirectURI  string
	Scope        string
	ResponseType string
	AccessType   string
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
		ClientID:    props.ClientID,
		RedirectURL: props.RedirectURI,
		Scopes:      []string{props.Scope},
		Endpoint:    google.Endpoint,
	}
}
