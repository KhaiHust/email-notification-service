package controller

import "github.com/KhaiHust/email-notification-service/core/helper"

type BaseController struct {
	validator *helper.CustomValidate
}

func NewBaseController(validator *helper.CustomValidate) *BaseController {
	return &BaseController{
		validator: validator,
	}
}
