package properties

import "github.com/golibs-starter/golib/config"

type NewrelicProperties struct {
	AppName    string
	LicenseKey string
	Enabled    bool
}

func (n NewrelicProperties) Prefix() string {
	return "app.services.newrelic"
}

func NewNewrelicProperties(loader config.Loader) (*NewrelicProperties, error) {
	props := &NewrelicProperties{}
	if err := loader.Bind(props); err != nil {
		return nil, err
	}
	return props, nil
}
