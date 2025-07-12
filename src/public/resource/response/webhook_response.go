package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type WebhookResponse struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
}

func ToWebhookResponse(webEntity *entity.WebhookEntity) *WebhookResponse {
	if webEntity == nil {
		return nil
	}
	return &WebhookResponse{
		ID:      webEntity.ID,
		Type:    webEntity.Type,
		Enabled: webEntity.Enabled,
		Name:    webEntity.Name,
		URL:     webEntity.URL,
	}
}
