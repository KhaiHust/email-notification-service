package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailRequestResponse struct {
	ID          int64  `json:"id"`
	WorkspaceID int64  `json:"workspace_id"`
	TemplateID  int64  `json:"template_id"`
	Status      string `json:"status"`
	SentAt      *int64 `json:"sent_at,omitempty"`
	RequestID   string `json:"request_id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func ToEmailRequestResponse(emailRqEntity *entity.EmailRequestEntity) *EmailRequestResponse {
	if emailRqEntity == nil {
		return nil
	}
	return &EmailRequestResponse{
		ID:          emailRqEntity.ID,
		WorkspaceID: emailRqEntity.WorkspaceID,
		TemplateID:  emailRqEntity.TemplateId,
		Status:      emailRqEntity.Status,
		SentAt:      emailRqEntity.SentAt,
		RequestID:   emailRqEntity.RequestID,
		CreatedAt:   emailRqEntity.CreatedAt,
		UpdatedAt:   emailRqEntity.UpdatedAt,
	}
}
func ToListEmailRequestResponse(emailRqEntities []*entity.EmailRequestEntity) []*EmailRequestResponse {
	if emailRqEntities == nil {
		return nil
	}
	emailRqResponses := make([]*EmailRequestResponse, len(emailRqEntities))
	for i, emailRqEntity := range emailRqEntities {
		emailRqResponses[i] = ToEmailRequestResponse(emailRqEntity)
	}
	return emailRqResponses
}
