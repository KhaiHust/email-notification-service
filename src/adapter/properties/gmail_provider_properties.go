package properties

import "github.com/golibs-starter/golib/config"

type GmailProviderProperties struct {
	BaseURL      string
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
