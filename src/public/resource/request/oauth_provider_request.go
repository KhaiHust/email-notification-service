package request

type CreateEmailProviderRequest struct {
	Code string `json:"code" validate:"required"`
}
