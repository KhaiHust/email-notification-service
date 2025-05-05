package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICreateApiKeyUseCase interface {
	GenerateApiKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
}
type CreateApiKeyUseCase struct {
	encryptUseCase       IEncryptUseCase
	apiKeyRepositoryPort port.IApiKeyRepositoryPort
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
	hash := sha256.Sum256([]byte(encrypted))
	apiKeyEntity.KeyHash = hex.EncodeToString(hash[:])
	apiKeyEntity.RawPrefix = prefix

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
) ICreateApiKeyUseCase {
	return &CreateApiKeyUseCase{
		encryptUseCase:       encryptUseCase,
		apiKeyRepositoryPort: apiKeyRepositoryPort,
	}
}
