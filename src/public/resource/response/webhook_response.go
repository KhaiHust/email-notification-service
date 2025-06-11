package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type WebhookResponse struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func ToWebhookResponse(webEntity *entity.WebhookEntity) *WebhookResponse {
	if webEntity == nil {
		return nil
	}
	return &WebhookResponse{
		ID:      webEntity.ID,
		Type:    webEntity.Type,
		Enabled: webEntity.Enabled,
	}
}
