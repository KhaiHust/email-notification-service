package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetWorkspaceUseCase interface {
	GetWorkspaceByCode(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error)
}
type GetWorkspaceUseCase struct {
	workspaceRepositoryPort port.IWorkspaceRepositoryPort
}

func (g GetWorkspaceUseCase) GetWorkspaceByCode(ctx context.Context, userId int64, code string) (*entity.WorkspaceEntity, error) {
	return g.workspaceRepositoryPort.GetWorkspaceByCodeAndUserId(ctx, userId, code)
}

func NewGetWorkspaceUseCase(workspaceRepositoryPort port.IWorkspaceRepositoryPort) IGetWorkspaceUseCase {
	return &GetWorkspaceUseCase{
		workspaceRepositoryPort: workspaceRepositoryPort,
	}
}
