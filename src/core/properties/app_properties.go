package properties

import "github.com/golibs-starter/golib/config"

type AppProperties struct {
	GoogleCloudProject string `default:""`
}

func (a AppProperties) Prefix() string {
	return "app.properties"
}

func NewAppProperties(loader config.Loader) (*AppProperties, error) {
	props := AppProperties{}
	err := loader.Bind(&props)
	return &props, err
}
