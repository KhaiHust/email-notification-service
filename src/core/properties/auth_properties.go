package properties

import "github.com/golibs-starter/golib/config"

type AuthProperties struct {
	ExpiredAccessToken  int64
	ExpiredRefreshToken int64
	PrivateKey          string
}

func (a AuthProperties) Prefix() string {
	return "app.services.auth.jwt"
}

func NewAuthProperties(loader config.Loader) (*AuthProperties, error) {
	props := AuthProperties{}
	err := loader.Bind(&props)
	return &props, err
}
