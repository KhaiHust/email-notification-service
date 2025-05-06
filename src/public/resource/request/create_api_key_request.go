package request

type CreateApiKeyRequest struct {
	Name        string `json:"name" validate:"required"`
	Environment string `json:"environment" validate:"required"`
}
