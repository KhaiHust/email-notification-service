package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type WorkspaceUserRepositoryAdapter struct {
	base
}

func (w WorkspaceUserRepositoryAdapter) GetWorkspaceUserByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WorkspaceUserEntity, error) {
	var workspaceUserModels []*model.WorkspaceUserModel
	err := w.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&workspaceUserModels).Error
	if err != nil {
		return nil, err
	}
	return mapper.ToListWorkspaceUserEntity(workspaceUserModels), nil
}

func (w WorkspaceUserRepositoryAdapter) GetWorkspaceUserByWorkspaceIDAndUserID(ctx context.Context, workspaceID int64, userID int64) (*entity.WorkspaceUserEntity, error) {
	var workspaceUserModel model.WorkspaceUserModel
	err := w.db.WithContext(ctx).Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&workspaceUserModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToWorkspaceUserEntity(&workspaceUserModel), nil
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
