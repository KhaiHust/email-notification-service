package request

import "github.com/KhaiHust/email-notification-service/core/entity"

type CreateEmailTemplateRequest struct {
	Name    string `json:"name" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Version string `json:"version" `
}

func ToEmailTemplateEntity(req *CreateEmailTemplateRequest) *entity.EmailTemplateEntity {
	return &entity.EmailTemplateEntity{
		Name:    req.Name,
		Subject: req.Subject,
		Body:    req.Body,
		Version: req.Version,
	}
}
