package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

const (
	PrefixEmailProviderKeyCache = "email_provider_cache:id:%d"
)

type IGetEmailProviderUseCase interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
	GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID int64, workspaceID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*entity.EmailProviderEntity, error)
	GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error)
	GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error)
	GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx context.Context, provider string, workspaceID int64, environment string) (*entity.EmailProviderEntity, error)
}
type GetEmailProviderUseCase struct {
	emailProviderPort           port.IEmailProviderPort
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
	encryptUseCase              IEncryptUseCase
	redisPort                   port.IRedisPort
}

func (g GetEmailProviderUseCase) GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx context.Context, provider string, workspaceID int64, environment string) (*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx, provider, workspaceID, environment)
}

func (g GetEmailProviderUseCase) GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error) {
	providers, err := g.emailProviderRepositoryPort.GetProvidersByIds(ctx, ids)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("GetProvidersByIds error: %v", err))
		return nil, err
	}
	for _, provider := range providers {
		if provider.OAuthToken != "" {
			decryptedToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthToken)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("DecryptProviderToken error: %v", err))
				return nil, err
			}
			provider.OAuthToken = decryptedToken
		}
		if provider.OAuthRefreshToken != "" {
			decryptedRefreshToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthRefreshToken)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("DecryptProviderToken error: %v", err))
				return nil, err
			}
			provider.OAuthRefreshToken = decryptedRefreshToken
		}
	}
	return providers, nil
}

func (g GetEmailProviderUseCase) GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetAllEmailProviders(ctx, filter)
}

func (g GetEmailProviderUseCase) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetEmailProviderByWorkspaceCodeAndProvider(ctx, workspaceCode, provider)
}

func (g GetEmailProviderUseCase) GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID, workspaceID int64) (*entity.EmailProviderEntity, error) {

	provider, err := g.emailProviderRepositoryPort.GetEmailProviderByIDAndWorkspaceID(ctx, providerID, workspaceID)
	if err != nil {
		log.Error(ctx, "GetEmailProviderByIDAndWorkspaceID error: %v", err)
		return nil, err
	}
	//decrypt OAuth tokens if they exist
	err = g.decryptToken(ctx, provider)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (g GetEmailProviderUseCase) decryptToken(ctx context.Context, provider *entity.EmailProviderEntity) error {
	if provider.OAuthToken != "" {
		decryptedToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthToken)
		if err != nil {
			log.Error(ctx, "DecryptProviderToken error: %v", err)
			return err
		}
		provider.OAuthToken = decryptedToken
	}
	if provider.OAuthRefreshToken != "" {
		decryptedRefreshToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthRefreshToken)
		if err != nil {
			log.Error(ctx, "DecryptProviderToken error: %v", err)
			return err
		}
		provider.OAuthRefreshToken = decryptedRefreshToken
	}
	return nil
}

func (g GetEmailProviderUseCase) GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error) {
	// Check Redis cache first
	cacheKey := fmt.Sprintf(PrefixEmailProviderKeyCache, ID)
	cachedData, err := g.redisPort.GetFromRedis(ctx, cacheKey)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error getting data from Redis: %v", err))
		return nil, err
	}
	if cachedData != nil {
		var emailProvider entity.EmailProviderEntity
		if err := json.Unmarshal(cachedData, &emailProvider); err != nil {
			log.Error(ctx, fmt.Sprintf("Error unmarshalling cached data: %v", err))
			return nil, err
		}
		err = g.decryptToken(ctx, &emailProvider)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Error decrypting token: %v", err))
			return nil, err
		}
		return &emailProvider, nil
	}
	emailProvider, err := g.emailProviderRepositoryPort.GetEmailProviderByID(ctx, ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("GetEmailProviderByID error: %v", err))
		return nil, err
	}
	if err = g.decryptToken(ctx, emailProvider); err != nil {
		log.Error(ctx, fmt.Sprintf("Error decrypting token: %v", err))
		return nil, err
	}
	go func() {
		if err := g.redisPort.SetToRedis(ctx, cacheKey, &emailProvider, 3600); err != nil { // Cache for 1 hour
			log.Error(ctx, fmt.Sprintf("Error setting data to Redis: %v", err))
		}
	}()
	return emailProvider, nil
}

func (g GetEmailProviderUseCase) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error) {
	return g.emailProviderPort.GetOAuthUrl(ctx, provider)
}

func NewGetEmailProviderUseCase(
	emailProviderPort port.IEmailProviderPort,
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
	encryptUseCase IEncryptUseCase,
	redisPort port.IRedisPort,
) IGetEmailProviderUseCase {
	return &GetEmailProviderUseCase{
		emailProviderPort:           emailProviderPort,
		emailProviderRepositoryPort: emailProviderRepositoryPort,
		encryptUseCase:              encryptUseCase,
		redisPort:                   redisPort,
	}
}
