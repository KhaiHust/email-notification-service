package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type WorkspaceResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func ToWorkspaceResponse(wEntity *entity.WorkspaceEntity) *WorkspaceResponse {
	return &WorkspaceResponse{
		ID:          wEntity.ID,
		Name:        wEntity.Name,
		Code:        wEntity.Code,
		Description: wEntity.Description,
	}
}
