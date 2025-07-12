package properties

type CORSProperties struct {
	AllowOrigins     []string
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	AllowCredentials bool
	MaxAge           int
	AllowOriginFunc  string
}
