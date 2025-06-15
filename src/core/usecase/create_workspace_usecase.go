package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/google/uuid"
)

type ICreateWorkspaceUseCase interface {
	CreateNewWorkspace(ctx context.Context, userID int64, workspaceEntity *entity.WorkspaceEntity) (*entity.WorkspaceEntity, error)
}
type CreateWorkspaceUseCase struct {
	workspaceRepo              port.IWorkspaceRepositoryPort
	workspaceUserRepoPort      port.IWorkspaceUserRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
	createApiKeyUseCase        ICreateApiKeyUseCase
}

func (c CreateWorkspaceUseCase) CreateNewWorkspace(ctx context.Context, userID int64, workspaceEntity *entity.WorkspaceEntity) (*entity.WorkspaceEntity, error) {

	workspaceCode := uuid.New().String()
	workspaceEntity.Code = workspaceCode

	tx := c.databaseTransactionUseCase.StartTx()
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
			log.Error(ctx, "CreateNewWorkspace", "Panic when create new workspace: %v", r)
		}
		if err != nil {
			if err = c.databaseTransactionUseCase.RollbackTx(tx); err != nil {
				log.Error(ctx, "CreateNewWorkspace", "Error when rollback tx: %v", err)
			} else {
				log.Info(ctx, "CreateNewWorkspace", "Rollback tx successfully")
			}
		}
	}()
	workspace, err := c.workspaceRepo.SaveWorkspace(ctx, tx, workspaceEntity)
	if err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when save workspace: %v", err)
		return nil, err
	}
	workspaceUser := &entity.WorkspaceUserEntity{
		UserID:      userID,
		WorkspaceID: workspace.ID,
		Role:        constant.WorkspaceRoleAdmin,
	}
	_, err = c.workspaceUserRepoPort.SaveNewWorkspaceUser(ctx, tx, workspaceUser)
	if err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when save workspace user: %v", err)
		return nil, err
	}
	//create api key for product and test environment
	if err = c.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when commit tx: %v", err)
		return nil, err
	}
	return workspace, nil

}

func NewCreateWorkspaceUseCase(
	workspaceRepo port.IWorkspaceRepositoryPort,
	workspaceUserRepoPort port.IWorkspaceUserRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	createApiKeyUseCase ICreateApiKeyUseCase,
) ICreateWorkspaceUseCase {
	return &CreateWorkspaceUseCase{
		workspaceRepo:              workspaceRepo,
		workspaceUserRepoPort:      workspaceUserRepoPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
		createApiKeyUseCase:        createApiKeyUseCase,
	}
}
