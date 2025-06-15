package request

import "github.com/KhaiHust/email-notification-service/core/entity/dto/request"

type CreateWebhookRequest struct {
	URL     string `json:"url" validate:"required,url"`
	Type    string `json:"type" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Enabled bool   `json:"enabled" default:"true"`
}

func ToCreateWebhookRequestDto(req *CreateWebhookRequest) *request.CreateWebhookRequest {
	if req == nil {
		return nil
	}
	return &request.CreateWebhookRequest{
		URL:     req.URL,
		Type:    req.Type,
		Name:    req.Name,
		Enabled: req.Enabled,
	}
}
