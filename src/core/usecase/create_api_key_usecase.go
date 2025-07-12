package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICreateApiKeyUseCase interface {
	CreateNewApiKey(ctx context.Context, apiKey *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
	RevokeAPIKeyByID(ctx context.Context, workspaceID, apiKeyID int64, req *request.RevokeApiKeyRequest) error
}
type CreateApiKeyUseCase struct {
	encryptUseCase             IEncryptUseCase
	apiKeyRepositoryPort       port.IApiKeyRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
}

func (c CreateApiKeyUseCase) RevokeAPIKeyByID(ctx context.Context, workspaceID, apiKeyID int64, req *request.RevokeApiKeyRequest) error {
	// get API key by ID and workspace ID
	apiKeyEntity, err := c.apiKeyRepositoryPort.GetAPIKeyByIDAndWorkspaceID(ctx, apiKeyID, workspaceID)
	if err != nil {
		log.Error(ctx, "Error getting API key by ID and workspace ID: %v", err)
		return err
	}
	apiKeyEntity.Revoked = true
	tx := c.databaseTransactionUseCase.StartTx()
	defer func() {

		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := c.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Rollback error: %v", errRollback)
			} else {
				log.Info(ctx, "Rollback successfully")
			}
		}
	}()
	apiKeyEntity, err = c.apiKeyRepositoryPort.UpdateAPIKey(ctx, tx, apiKeyEntity)
	if err != nil {
		log.Error(ctx, "Error updating API key: %v", err)
		return err
	}
	if err := c.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "CommitTx error: %v", err)
		return err
	}
	return nil
}

func (c CreateApiKeyUseCase) CreateNewApiKey(ctx context.Context, apiKey *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error) {

	tx := c.databaseTransactionUseCase.StartTx()
	var err error
	defer func() {

		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := c.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Rollback error: %v", errRollback)
			} else {
				log.Info(ctx, "Rollback successfully")
			}
		}
	}()
	apiKeyEntity, err := c.GenerateApiKey(ctx, tx, apiKey)
	if err != nil {
		log.Error(ctx, "GenerateApiKey error: %v", err)
		return nil, err
	}
	if err := c.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "CommitTx error: %v", err)
		return nil, err
	}
	return apiKeyEntity, nil
}

func (c CreateApiKeyUseCase) GenerateApiKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error) {

	prefix := uuid.New().String()[:12]
	secret := uuid.New().String()

	rawKey := fmt.Sprintf("%d_%s.%s", apiKeyEntity.WorkspaceID, prefix, secret)

	hash := sha256.Sum256([]byte(rawKey))
	apiKeyEntity.KeyHash = hex.EncodeToString(hash[:])
	apiKeyEntity.RawPrefix = prefix

	apiKeyEntity, err := c.apiKeyRepositoryPort.SaveNewApiKey(ctx, tx, apiKeyEntity)
	if err != nil {
		log.Error(ctx, "Error saving new api key: %v", err)
		return nil, err
	}
	apiKeyEntity.RawKey = rawKey

	return apiKeyEntity, nil
}

func NewCreateApiKeyUseCase(
	encryptUseCase IEncryptUseCase,
	apiKeyRepositoryPort port.IApiKeyRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
) ICreateApiKeyUseCase {
	return &CreateApiKeyUseCase{
		encryptUseCase:             encryptUseCase,
		apiKeyRepositoryPort:       apiKeyRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
	}
}
