package entity

type EmailLogsEntity struct {
	BaseEntity
	EmailRequestID  int64
	TemplateId      int64
	Recipient       string
	Status          string
	ErrorMessage    string
	RetryCount      int64
	RequestID       string
	WorkspaceID     int64
	EmailProviderID int64
	LoggedAt        int64
}
