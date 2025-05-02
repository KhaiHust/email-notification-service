package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IWorkspaceService interface {
	CreateNewWorkspace(ctx context.Context, userID int64, req *request.CreateWorkspaceRequest) (*response.WorkspaceResponse, error)
	GetWorkspacesByUserId(ctx context.Context, userID int64) ([]*response.WorkspaceResponse, error)
}
type WorkspaceService struct {
	createWorkspaceUseCase usecase.ICreateWorkspaceUseCase
	getWorkspaceUseCase    usecase.IGetWorkspaceUseCase
}

func (w WorkspaceService) GetWorkspacesByUserId(ctx context.Context, userID int64) ([]*response.WorkspaceResponse, error) {
	workspaceEntities, err := w.getWorkspaceUseCase.GetWorkspaceByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.ToListWorkspaceResponse(workspaceEntities), nil
}

func (w WorkspaceService) CreateNewWorkspace(ctx context.Context, userID int64, req *request.CreateWorkspaceRequest) (*response.WorkspaceResponse, error) {
	workspaceEntity := request.ToCreateWorkspaceEntity(req)
	workspace, err := w.createWorkspaceUseCase.CreateNewWorkspace(ctx, userID, workspaceEntity)
	if err != nil {
		return nil, err
	}
	return response.ToWorkspaceResponse(workspace), nil
}

func NewWorkspaceService(
	createWorkspaceUseCase usecase.ICreateWorkspaceUseCase,
	getWorkspaceUseCase usecase.IGetWorkspaceUseCase,
) IWorkspaceService {
	return &WorkspaceService{
		createWorkspaceUseCase: createWorkspaceUseCase,
		getWorkspaceUseCase:    getWorkspaceUseCase,
	}
}
