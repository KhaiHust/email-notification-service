package properties

import "github.com/golibs-starter/golib/config"

type TaskProperties struct {
	BaseUrl   string
	SecretKey string
}

func (t TaskProperties) Prefix() string {
	return "app.services.task"
}

func NewTaskProperties(loader config.Loader) (*TaskProperties, error) {
	var taskProperties TaskProperties
	if err := loader.Bind(&taskProperties); err != nil {
		return nil, err
	}
	return &taskProperties, nil
}
