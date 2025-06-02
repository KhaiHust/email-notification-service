package request

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailDataDto struct {
	EmailRequestID int64
	TrackingID     string
	To             string
	Subject        string
	Body           string
	Provider       *entity.EmailProviderEntity
}
