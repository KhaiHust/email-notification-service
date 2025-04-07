package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

type IWorkspaceRepositoryPort interface {
	GetWorkspaceByCodeAndUserId(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error)
}
