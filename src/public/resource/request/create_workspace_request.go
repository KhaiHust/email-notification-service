package request

import "github.com/KhaiHust/email-notification-service/core/entity"

type CreateWorkspaceRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func ToCreateWorkspaceEntity(req *CreateWorkspaceRequest) *entity.WorkspaceEntity {
	return &entity.WorkspaceEntity{
		Name:        req.Name,
		Description: req.Description,
	}
}
