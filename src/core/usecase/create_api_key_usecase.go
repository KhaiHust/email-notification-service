package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICreateApiKeyUseCase interface {
	GenerateApiKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
	CreateNewApiKey(ctx context.Context, apiKey *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
}
type CreateApiKeyUseCase struct {
	encryptUseCase             IEncryptUseCase
	apiKeyRepositoryPort       port.IApiKeyRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
}

func (c CreateApiKeyUseCase) CreateNewApiKey(ctx context.Context, apiKey *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error) {

	tx := c.databaseTransactionUseCase.StartTx()
	defer func() {
		var err error
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
	raw := fmt.Sprintf("%d_%s_%s.%s", apiKeyEntity.WorkspaceID, prefix, secret)
	encrypted, err := c.encryptUseCase.EncryptAES(ctx, raw)
	if err != nil {
		log.Error(ctx, "Error encrypting api key: %v", err)
		return nil, err
	}
	encryptedKey, err := c.encryptUseCase.EncryptAES(ctx, encrypted)
	if err != nil {
		log.Error(ctx, "Error encrypting api key: %v", err)
		return nil, err
	}
	hash := sha256.Sum256([]byte(encrypted))
	apiKeyEntity.KeyHash = hex.EncodeToString(hash[:])
	apiKeyEntity.RawPrefix = prefix
	apiKeyEntity.KeyEnc = encryptedKey
	apiKeyEntity, err = c.apiKeyRepositoryPort.SaveNewApiKey(ctx, tx, apiKeyEntity)
	if err != nil {
		log.Error(ctx, "Error saving new api key: %v", err)
		return nil, err
	}
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
