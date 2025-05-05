package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		if err = c.databaseTransactionUseCase.RollbackTx(tx); err != nil {
			log.Error(ctx, "CreateNewWorkspace", "Error when rollback tx: %v", err)
		} else {
			log.Info(ctx, "CreateNewWorkspace", "Rollback tx successfully")
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
	err = c.generateAPIKeyForWorkspace(ctx, workspace, tx)
	if err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when generate api key for workspace: %v", err)
		return nil, err
	}
	if err = c.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when commit tx: %v", err)
		return nil, err
	}
	return workspace, nil

}

func (c CreateWorkspaceUseCase) generateAPIKeyForWorkspace(ctx context.Context, workspace *entity.WorkspaceEntity, tx *gorm.DB) error {
	apiKeyEntity := &entity.ApiKeyEntity{
		WorkspaceID: workspace.ID,
		Name:        fmt.Sprintf("%s %s %s", "APIKey", workspace.Name, constant.EnvironmentProduction),
		Environment: constant.EnvironmentProduction,
	}
	_, err := c.createApiKeyUseCase.GenerateApiKey(ctx, tx, apiKeyEntity)
	if err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when create api key for production: %v", err)
		return err
	}
	apiKeyEntity.Name = fmt.Sprintf("%s %s %s", "APIKey", workspace.Name, constant.EnvironmentTest)
	apiKeyEntity.Environment = constant.EnvironmentTest
	_, err = c.createApiKeyUseCase.GenerateApiKey(ctx, tx, apiKeyEntity)
	if err != nil {
		log.Error(ctx, "CreateNewWorkspace", "Error when create api key for test: %v", err)
		return err
	}
	return nil
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
