package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
	"strconv"
)

type EmailLogController struct {
	BaseController  *BaseController
	emailLogService service.IEmailLogService
}

func (e *EmailLogController) GetLogs(ctx *gin.Context) {
	emailRequestIDStr := ctx.Param(constant.ParamEmailRequestId)
	if emailRequestIDStr == "" {
		log.Error(ctx, "Error when get email request id")
		apihelper.AbortErrorHandle(ctx, common.ErrBadRequest)
		return
	}
	emailRequestID, err := strconv.ParseInt(emailRequestIDStr, 10, 64)
	if err != nil {
		log.Error(ctx, "Error when parse email request id", err)
		apihelper.AbortErrorHandle(ctx, common.ErrBadRequest)
		return
	}
	workspaceID := e.BaseController.GetWorkspaceIDFromContext(ctx)
	if workspaceID == 0 {
		log.Error(ctx, "Error when get workspace id from context", common.ErrBadRequest)
		apihelper.AbortErrorHandle(ctx, common.ErrBadRequest)
		return
	}
	emailLogs, err := e.emailLogService.GetLogsByEmailRequestIDAndWorkspaceID(ctx, emailRequestID, workspaceID)
	if err != nil {
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	apihelper.SuccessfulHandle(ctx, response.ToListEmailLogResource(emailLogs))
}
func NewEmailLogController(
	base *BaseController,
	emailLogService service.IEmailLogService,
) *EmailLogController {
	return &EmailLogController{
		BaseController:  base,
		emailLogService: emailLogService,
	}
}
