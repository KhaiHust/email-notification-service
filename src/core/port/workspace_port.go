package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IWorkspaceRepositoryPort interface {
	GetWorkspaceByCodeAndUserId(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error)
	SaveWorkspace(ctx context.Context, db *gorm.DB, workspaceEntity *entity.WorkspaceEntity) (*entity.WorkspaceEntity, error)
	GetWorkspaceByUserId(ctx context.Context, userId int64) ([]*entity.WorkspaceEntity, error)
}
