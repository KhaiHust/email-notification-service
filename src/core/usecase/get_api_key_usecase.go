package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

const (
	PrefixCacheApiKey = "api_key:%s" // %s is the hash key of the API key
)

type IGetApiKeyUseCase interface {
	GetAll(ctx context.Context, filter *request.GetApiKeyRequestFilter) ([]*entity.ApiKeyEntity, error)
	GetApiKeyDetail(ctx context.Context, apiKey string) (*entity.ApiKeyEntity, error)
}
type GetApiKeyUseCase struct {
	redisPort            port.IRedisPort
	apiKeyRepositoryPort port.IApiKeyRepositoryPort
}

func (g GetApiKeyUseCase) GetApiKeyDetail(ctx context.Context, apiKey string) (*entity.ApiKeyEntity, error) {
	hash := sha256.Sum256([]byte(apiKey))
	hashEncoded := hex.EncodeToString(hash[:])
	// Check cache first
	cacheKey := fmt.Sprintf(PrefixCacheApiKey, hashEncoded)
	cachedAPIKey, err := g.redisPort.GetFromRedis(ctx, cacheKey)
	if err != nil {
		log.Error(ctx, "Error getting API key from cache: %v", err)
		return nil, err
	}
	apiKeyEntity := &entity.ApiKeyEntity{}
	if cachedAPIKey != nil {
		// parse the cached API key entity
		if err = json.Unmarshal(cachedAPIKey, apiKeyEntity); err != nil {
			log.Error(ctx, "Error unmarshalling cached API key: %v", err)
			return nil, err
		}

	} else {
		apiKeyEntity, err = g.apiKeyRepositoryPort.GetAPIKeyByHashKey(ctx, hashEncoded)
		if err != nil {
			log.Error(ctx, "Error getting API key from repository: %v", err)
			return nil, err
		}
		// Cache the API key entity
		go func() {
			if err = g.redisPort.SetToRedis(ctx, cacheKey, apiKeyEntity, constant.DefaultTTL); err != nil {
				log.Error(ctx, "Error setting API key to cache: %v", err)
			}
		}()
	}
	return apiKeyEntity, nil
}

func (g GetApiKeyUseCase) GetAll(ctx context.Context, filter *request.GetApiKeyRequestFilter) ([]*entity.ApiKeyEntity, error) {
	return g.apiKeyRepositoryPort.GetAllApiKeys(ctx, filter)
}

func NewGetApiKeyUseCase(
	redisPort port.IRedisPort,
	apiKeyRepositoryPort port.IApiKeyRepositoryPort,
) IGetApiKeyUseCase {
	return &GetApiKeyUseCase{
		redisPort:            redisPort,
		apiKeyRepositoryPort: apiKeyRepositoryPort,
	}
}
