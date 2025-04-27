package controller

import (
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailSendingController struct {
	*BaseController
	emailSendingService service.IEmailSendingService
}

func (e *EmailSendingController) SendEmailRequest(ctx *gin.Context) {
	//workspaceCode := ctx.Param(constant.ParamWorkspaceCode)
	var req request.EmailSendingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctx, "Error when binding request", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	err := e.emailSendingService.SendEmailRequest(ctx, 3, &req)
	if err != nil {
		log.Error(ctx, "Error when sending email", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	apihelper.SuccessfulHandle(ctx, nil)
}
func NewEmailSendingController(
	base *BaseController,
	emailSendingService service.IEmailSendingService,
) *EmailSendingController {
	return &EmailSendingController{
		BaseController:      base,
		emailSendingService: emailSendingService,
	}
}
