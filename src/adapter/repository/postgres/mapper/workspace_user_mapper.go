package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToWorkspaceUserModel(workspaceEntity *entity.WorkspaceUserEntity) *model.WorkspaceUserModel {
	return &model.WorkspaceUserModel{
		BaseModel: model.BaseModel{
			ID: workspaceEntity.ID,
		},
		WorkspaceID: workspaceEntity.WorkspaceID,
		UserID:      workspaceEntity.UserID,
		Role:        workspaceEntity.Role,
	}
}
func ToWorkspaceUserEntity(workspaceUserModel *model.WorkspaceUserModel) *entity.WorkspaceUserEntity {
	return &entity.WorkspaceUserEntity{
		BaseEntity: entity.BaseEntity{
			ID:        workspaceUserModel.ID,
			UpdatedAt: workspaceUserModel.UpdatedAt.Unix(),
			CreatedAt: workspaceUserModel.CreatedAt.Unix(),
		},
		WorkspaceID: workspaceUserModel.WorkspaceID,
		UserID:      workspaceUserModel.UserID,
		Role:        workspaceUserModel.Role,
	}
}
func ToListWorkspaceUserEntity(workspaceUserModels []*model.WorkspaceUserModel) []*entity.WorkspaceUserEntity {
	var workspaceUserEntities []*entity.WorkspaceUserEntity
	for _, workspaceUserModel := range workspaceUserModels {
		workspaceUserEntities = append(workspaceUserEntities, ToWorkspaceUserEntity(workspaceUserModel))
	}
	return workspaceUserEntities
}
