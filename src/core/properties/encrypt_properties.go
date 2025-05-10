package properties

import "github.com/golibs-starter/golib/config"

type EncryptProperties struct {
	EncryptKey             string
	EncryptVersionTemplate string
}

func (e EncryptProperties) Prefix() string {
	return "app.services.encrypt"
}

func NewEncryptProperties(loader config.Loader) (*EncryptProperties, error) {
	props := EncryptProperties{}
	err := loader.Bind(&props)
	return &props, err
}
