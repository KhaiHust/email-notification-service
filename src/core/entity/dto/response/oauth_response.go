package response

type OAuthUrlResponseDto struct {
	Url string
}
type OAuthInfoResponseDto struct {
	AccessToken  string
	RefreshToken string
	ExpiredAt    int64
	TokenType    string
	Scope        string
	Email        string
	SmtpHost     string
	SmtpPort     int
	UseTLS       bool
}
