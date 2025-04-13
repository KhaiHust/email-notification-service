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

type WorkspaceRepositoryAdapter struct {
	base
}

func (w WorkspaceRepositoryAdapter) SaveWorkspace(ctx context.Context, db *gorm.DB, workspaceEntity *entity.WorkspaceEntity) (*entity.WorkspaceEntity, error) {
	workspaceModel := mapper.ToWorkspaceModel(workspaceEntity)
	err := db.WithContext(ctx).Table(workspaceModel.TableName()).Create(workspaceModel).Error
	if err != nil {
		return nil, err
	}
	return mapper.ToWorkspaceEntity(workspaceModel), nil
}

func (w WorkspaceRepositoryAdapter) GetWorkspaceByCodeAndUserId(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error) {
	var workspaceModel *model.WorkspaceModel
	err := w.db.WithContext(ctx).Table(model.WorkspaceModel{}.TableName()+" AS w").
		Joins("left join workspace_users wu on wu.workspace_id = w.id").
		Where("w.code = ? AND wu.user_id = ?", code, userId).First(&workspaceModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToWorkspaceEntity(workspaceModel), nil
}

func NewWorkspaceRepositoryAdapter(db *gorm.DB) port.IWorkspaceRepositoryPort {
	return &WorkspaceRepositoryAdapter{
		base: base{db},
	}
}
