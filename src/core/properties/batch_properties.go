package properties

import "github.com/golibs-starter/golib/config"

type BatchProperties struct {
	BatchSize    int `default:"100"`
	NumOfWorkers int `default:"5"`
}

func (b BatchProperties) Prefix() string {
	return "app.configs.batch"
}

func NewBatchProperties(loader config.Loader) (*BatchProperties, error) {
	props := BatchProperties{}
	if err := loader.Bind(&props); err != nil {
		return nil, err
	}
	return &props, nil
}
