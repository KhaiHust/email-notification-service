package properties

import "github.com/golibs-starter/golib/config"

type TrackingProperties struct {
	BaseUrl string
}

func (t TrackingProperties) Prefix() string {
	return "app.services.tracking"
}

func NewTrackingProperties(loader config.Loader) (*TrackingProperties, error) {
	var props TrackingProperties
	err := loader.Bind(&props)
	if err != nil {
		return nil, err
	}
	return &props, nil
}
