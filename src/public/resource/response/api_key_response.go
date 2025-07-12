package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type ApiKeyResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	ExpiresAt   *int64 `json:"expires_at,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	RawKey      string `json:"raw_key,omitempty"`
}

func ToApiKeyResponse(key *entity.ApiKeyEntity) *ApiKeyResponse {
	if key == nil {
		return nil
	}
	return &ApiKeyResponse{
		ID:          key.ID,
		Name:        key.Name,
		Environment: key.Environment,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
		ExpiresAt:   key.ExpiresAt,
		Prefix:      key.RawPrefix,
		RawKey:      key.RawKey,
	}
}
func ToListApiKeyResponse(keys []*entity.ApiKeyEntity) []*ApiKeyResponse {
	if keys == nil {
		return nil
	}
	apiKeys := make([]*ApiKeyResponse, len(keys))
	for i, key := range keys {
		apiKeys[i] = ToApiKeyResponse(key)
	}
	return apiKeys
}
