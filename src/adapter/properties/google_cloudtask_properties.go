package properties

import "github.com/golibs-starter/golib/config"

type GoogleCloudTaskProperties struct {
	CloudTaskQueues  string `default:"queue1,queue2,queue3,queue4,queue5"`
	CloudTaskSaltKey string `default:""`
	GCPTaskLocation  string `default:"asia-east1"`
}

func (g GoogleCloudTaskProperties) Prefix() string {
	return "app.services.google-cloudtask"
}

func NewGoogleCloudTaskProperties(loader config.Loader) (*GoogleCloudTaskProperties, error) {
	props := GoogleCloudTaskProperties{}
	err := loader.Bind(&props)
	return &props, err
}
