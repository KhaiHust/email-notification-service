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
func ToListWorkspaceResponse(wEntities []*entity.WorkspaceEntity) []*WorkspaceResponse {
	workspaces := make([]*WorkspaceResponse, 0)
	for _, wEntity := range wEntities {
		workspaces = append(workspaces, ToWorkspaceResponse(wEntity))
	}
	return workspaces
}
