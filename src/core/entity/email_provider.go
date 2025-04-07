package entity

type EmailProviderEntity struct {
	BaseEntity
	WorkspaceId       int64
	Provider          string
	SmtpHost          string
	SmtpPort          int
	OAuthToken        string
	OAuthRefreshToken string
	OAuthExpiredAt    int64
	UseTLS            bool
}
