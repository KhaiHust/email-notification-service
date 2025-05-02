package request

import "github.com/KhaiHust/email-notification-service/core/entity/dto/request"

type CreateEmailProviderRequest struct {
	Code        string `json:"code" validate:"required"`
	FromName    string `json:"from_name" validate:"required"`
	Environment string `json:"environment" validate:"required,oneof=production test"`
}

func ToEmailProviderDto(req *CreateEmailProviderRequest) *request.CreateEmailProviderDto {
	return &request.CreateEmailProviderDto{
		Code:        req.Code,
		FromName:    req.FromName,
		Environment: req.Environment,
	}
}
