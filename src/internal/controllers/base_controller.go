package controllers

import (
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/helper"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
	validator *helper.CustomValidate
}

func (b *BaseController) GetWorkspaceIDFromContext(c *gin.Context) int64 {
	workspaceID, isExisted := c.Get(constant.WorkspaceIdKey)
	if !isExisted || workspaceID == "" {
		return 0
	}
	_, ok := workspaceID.(int64)
	if !ok {
		return 0
	}
	return workspaceID.(int64)
}
func NewBaseController(validator *helper.CustomValidate) *BaseController {
	return &BaseController{
		validator: validator,
	}
}
