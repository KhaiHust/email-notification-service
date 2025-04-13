package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IWorkspaceUserRepositoryPort interface {
	SaveNewWorkspaceUser(ctx context.Context, db *gorm.DB, workspaceUserEntity *entity.WorkspaceUserEntity) (*entity.WorkspaceUserEntity, error)
}
