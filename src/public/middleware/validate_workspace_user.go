package middleware

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib-security/web/context"
	"strconv"
)

type WorkspaceAccessMiddleware struct {
	validateAccessWorkspaceUsecase usecase.IValidateAccessWorkspaceUsecase
}

func NewWorkspaceAccessMiddleware(validateAccessWorkspaceUsecase usecase.IValidateAccessWorkspaceUsecase) *WorkspaceAccessMiddleware {
	return &WorkspaceAccessMiddleware{
		validateAccessWorkspaceUsecase: validateAccessWorkspaceUsecase,
	}
}

func (w *WorkspaceAccessMiddleware) WorkspaceAccessMiddlewareHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		workspaceCode := c.Param(constant.ParamWorkspaceCode)
		if workspaceCode == "" {
			apihelper.AbortErrorHandle(c, common.ErrForbidden)
			return
		}
		userIDStr := context.GetUserDetails(c.Request).Username()
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			apihelper.AbortErrorHandle(c, common.ErrForbidden)
			return
		}
		mapRole, workspaceId, err := w.validateAccessWorkspaceUsecase.ValidateAccessWorkspaceByUserIdAndCode(c, userID, workspaceCode)
		if err != nil {
			apihelper.AbortErrorHandle(c, common.ErrForbidden)
			return
		}
		if mapRole == nil {
			apihelper.AbortErrorHandle(c, common.ErrForbidden)
			return
		}
		c.Set(constant.WorkspaceRoleKey, mapRole[workspaceCode])
		c.Set(constant.WorkspaceIdKey, workspaceId)
	}
}
func ValidateRoleAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(constant.WorkspaceRoleKey)
		if !exists || role != constant.WorkspaceRoleAdmin {
			apihelper.AbortErrorHandle(c, common.ErrForbidden)
			return
		}
		c.Next()
	}
}
