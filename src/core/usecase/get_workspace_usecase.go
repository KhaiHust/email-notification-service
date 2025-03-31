package usecase

import "github.com/KhaiHust/email-notification-service/core/port"

type IGetWorkspaceUseCase interface {
}
type GetWorkspaceUseCase struct {
	workspaceRepositoryPort port.IWorkspaceRepositoryPort
}

func NewGetWorkspaceUseCase(workspaceRepositoryPort port.IWorkspaceRepositoryPort) IGetWorkspaceUseCase {
	return &GetWorkspaceUseCase{
		workspaceRepositoryPort: workspaceRepositoryPort,
	}
}
