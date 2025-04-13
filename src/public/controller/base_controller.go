package controller

import (
	"github.com/KhaiHust/email-notification-service/core/helper"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib-security/web/context"
	"strconv"
)

type BaseController struct {
	validator *helper.CustomValidate
}

func (b *BaseController) GetUserIDFromContext(c *gin.Context) (int64, error) {
	userIDStr := context.GetUserDetails(c.Request).Username()
	return strconv.ParseInt(userIDStr, 10, 64)
}
func NewBaseController(validator *helper.CustomValidate) *BaseController {
	return &BaseController{
		validator: validator,
	}
}
