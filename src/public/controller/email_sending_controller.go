package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	publicResponse "github.com/KhaiHust/email-notification-service/public/resource/response"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailSendingController struct {
	*BaseController
	emailSendingService service.IEmailSendingService
}

func (e *EmailSendingController) SendEmailRequest(ctx *gin.Context) {
	workspaceID := e.GetWorkspaceIDFromContext(ctx)
	if workspaceID == 0 {
		log.Error(ctx, "Error when get workspace id from context", common.ErrBadRequest)
		apihelper.AbortErrorHandle(ctx, common.ErrBadRequest)
		return
	}
	var req request.EmailSendingRequest
	if err := ctx.ShouldBind(&req); err != nil {
		log.Error(ctx, "Error when binding request", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	response, err := e.emailSendingService.SendEmailRequest(ctx, workspaceID, &req)
	if err != nil {
		log.Error(ctx, "Error when sending email", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	apihelper.SuccessfulHandle(ctx, &publicResponse.EmailSendingResponse{RequestID: response})
}

// SendEmailByTask handles the request to send an email by task ID
func (e *EmailSendingController) SendEmailByTask(ctx *gin.Context) {
	var req request.EmailSendTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error(ctx, "Error when binding request", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	if err := e.validator.Struct(&req); err != nil {
		log.Error(ctx, "Error when validating request", err)
		apihelper.AbortErrorHandle(ctx, common.ErrBadRequest)
		return
	}
	err := e.emailSendingService.SendEmailByTask(ctx, req.EmailRequestID)
	if err != nil {
		log.Error(ctx, "Error when sending email by task", err)
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
