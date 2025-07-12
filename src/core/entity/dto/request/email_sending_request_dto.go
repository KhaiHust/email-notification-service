package request

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailSendingRequestDto struct {
	TemplateID  int64
	WorkspaceID int64
	ProviderID  int64
	Datas       []*EmailSendingData
	Provider    *entity.EmailProviderEntity
}
type EmailSendingData struct {
	EmailRequestID int64
	TrackingID     string
	To             string
	SendAt         *int64
	IsRetry        bool
	Subject        map[string]string
	Body           map[string]string
}
