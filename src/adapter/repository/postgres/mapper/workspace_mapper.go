package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToWorkspaceModel(workspaceEntity *entity.WorkspaceEntity) *model.WorkspaceModel {
	return &model.WorkspaceModel{
		BaseModel: model.BaseModel{
			ID: workspaceEntity.ID,
		},
		Name:        workspaceEntity.Name,
		Description: workspaceEntity.Description,
		Code:        workspaceEntity.Code,
	}
}
func ToWorkspaceEntity(workspaceModel *model.WorkspaceModel) *entity.WorkspaceEntity {
	return &entity.WorkspaceEntity{
		BaseEntity: entity.BaseEntity{
			ID:        workspaceModel.ID,
			CreatedAt: workspaceModel.CreatedAt.Unix(),
			UpdatedAt: workspaceModel.UpdatedAt.Unix(),
		},
		Name:                workspaceModel.Name,
		Description:         workspaceModel.Description,
		Code:                workspaceModel.Code,
		WorkspaceUserEntity: ToListWorkspaceUserEntity(workspaceModel.WorkspaceUserModel),
	}
}
func ToListWorkspaceEntity(workspaceModels []*model.WorkspaceModel) []*entity.WorkspaceEntity {
	workspaces := make([]*entity.WorkspaceEntity, len(workspaceModels))
	for i, workspaceModel := range workspaceModels {
		workspaces[i] = ToWorkspaceEntity(workspaceModel)
	}
	return workspaces
}
