package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type WorkspaceUserRepositoryAdapter struct {
	base
}

func (w WorkspaceUserRepositoryAdapter) SaveNewWorkspaceUser(ctx context.Context, tx *gorm.DB, workspaceUserEntity *entity.WorkspaceUserEntity) (*entity.WorkspaceUserEntity, error) {
	workspaceUserModel := mapper.ToWorkspaceUserModel(workspaceUserEntity)
	err := tx.WithContext(ctx).Model(&model.WorkspaceUserModel{}).Create(workspaceUserModel).Error
	if err != nil {
		return nil, err
	}
	return workspaceUserEntity, nil
}

func NewWorkspaceUserRepositoryAdapter(db *gorm.DB) port.IWorkspaceUserRepositoryPort {
	return &WorkspaceUserRepositoryAdapter{
		base: base{db},
	}
}
