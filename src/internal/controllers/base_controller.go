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
func (b *BaseController) GetEnvironmentFromContext(c *gin.Context) string {
	environment, isExisted := c.Get(constant.EnvironmentKey)
	if !isExisted || environment == "" {
		return ""
	}
	env, ok := environment.(string)
	if !ok {
		return ""
	}
	return env
}
func NewBaseController(validator *helper.CustomValidate) *BaseController {
	return &BaseController{
		validator: validator,
	}
}
