package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToWorkspaceModel(entity *entity.WorkspaceEntity) *model.WorkspaceModel {
	return &model.WorkspaceModel{
		BaseModel: model.BaseModel{
			ID: entity.ID,
		},
		Name:        entity.Name,
		Description: entity.Description,
	}
}
func ToWorkspaceEntity(workspaceModel *model.WorkspaceModel) *entity.WorkspaceEntity {
	return &entity.WorkspaceEntity{
		BaseEntity: entity.BaseEntity{
			ID:        workspaceModel.ID,
			CreatedAt: workspaceModel.CreateAt.Unix(),
			UpdatedAt: workspaceModel.UpdateAt.Unix(),
		},
		Name:        workspaceModel.Name,
		Description: workspaceModel.Description,
	}
}
