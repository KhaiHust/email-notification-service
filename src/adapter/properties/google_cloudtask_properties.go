package properties

import "github.com/golibs-starter/golib/config"

type GoogleCloudTaskProperties struct {
	CloudTaskQueues  string `default:"queue01,queue02,queue03,queue04,queue05,queue06,queue07,queue08,queue09,queue10"`
	CloudTaskSaltKey string `default:""`
	GCPTaskLocation  string `default:"asia-east2"`
}

func (g GoogleCloudTaskProperties) Prefix() string {
	return "app.services.google-cloudtask"
}

func NewGoogleCloudTaskProperties(loader config.Loader) (*GoogleCloudTaskProperties, error) {
	props := GoogleCloudTaskProperties{}
	err := loader.Bind(&props)
	return &props, err
}
