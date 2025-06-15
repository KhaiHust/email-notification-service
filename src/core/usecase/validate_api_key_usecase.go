package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/golibs-starter/golib/log"
	"strconv"
	"strings"
	"time"
)

type IValidateApiKeyUsecase interface {
	ValidateKey(ctx context.Context, rawKey string) (bool, error)
}
type ValidateApiKeyUsecase struct {
	getApiKeyUseCase IGetApiKeyUseCase
}

func (v ValidateApiKeyUsecase) ValidateKey(ctx context.Context, rawKey string) (bool, error) {
	parts := strings.Split(rawKey, "_")
	if len(parts) < 2 {
		return false, nil
	}
	workspaceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.Error(ctx, "Error parsing workspace ID from rawKey: %v", err)
		return false, err
	}
	partKeys := strings.Split(parts[1], ".")
	if len(partKeys) < 2 {
		log.Error(ctx, "Invalid API key format: %s", rawKey)
		return false, nil
	}
	apiKeyEntity, err := v.getApiKeyUseCase.GetApiKeyDetail(ctx, rawKey)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error getting API key details for rawKey %s: %v", rawKey, err))
		return false, err
	}
	if apiKeyEntity == nil || apiKeyEntity.WorkspaceID != workspaceID ||
		apiKeyEntity.RawPrefix != partKeys[0] {
		log.Error(ctx, "API key validation failed for rawKey: %s", rawKey)
		return false, nil
	}
	if apiKeyEntity.Revoked || (apiKeyEntity.ExpiresAt != nil && *apiKeyEntity.ExpiresAt < time.Now().Unix()) {
		log.Error(ctx, "API key has expired for rawKey: %s", rawKey)
		return false, nil
	}
	ctx = context.WithValue(ctx, constant.WorkspaceIdKey, apiKeyEntity.WorkspaceID)
	ctx = context.WithValue(ctx, constant.EnvironmentKey, apiKeyEntity.Environment)
	return true, nil
}

func NewValidateApiKeyUsecase(
	getApiKeyUseCase IGetApiKeyUseCase,
) IValidateApiKeyUsecase {
	return &ValidateApiKeyUsecase{
		getApiKeyUseCase: getApiKeyUseCase,
	}
}
