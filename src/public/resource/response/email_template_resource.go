package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailTemplateResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	WorkspaceID int64  `json:"workspace_id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func ToEmailTemplateResponse(templateEntity *entity.EmailTemplateEntity) *EmailTemplateResponse {
	return &EmailTemplateResponse{
		Id:          templateEntity.ID,
		Name:        templateEntity.Name,
		WorkspaceID: templateEntity.WorkspaceId,
		CreatedAt:   templateEntity.CreatedAt,
		UpdatedAt:   templateEntity.UpdatedAt,
	}
}
