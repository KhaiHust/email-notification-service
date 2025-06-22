package controllers

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/internal/apihelper"
	"github.com/KhaiHust/email-notification-service/internal/resources/request"
	"github.com/KhaiHust/email-notification-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailSendingController struct {
	base                *BaseController
	emailSendingService services.IEmailSendingService
}

func (esc *EmailSendingController) SendEmailRequest(c *gin.Context) {
	workspaceID := esc.base.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is not provided in context")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	environment := esc.base.GetEnvironmentFromContext(c)
	if environment == "" ||
		(environment != constant.EnvironmentProduction &&
			environment != constant.EnvironmentTest) {
		log.Error(c, "Invalid environment provided in context")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	var req request.EmailSendingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "Failed to bind JSON request", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := esc.base.validator.Struct(req); err != nil {
		log.Error(c, "Validation failed for email sending request", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	req.Environment = environment
	resp, err := esc.emailSendingService.SendEmailRequest(c, workspaceID, &req)
	if err != nil {
		log.Error(c, "Failed to send email request", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, resp)
}
func NewEmailSendingController(
	base *BaseController,
	emailSendingService services.IEmailSendingService,
) *EmailSendingController {
	return &EmailSendingController{
		base:                base,
		emailSendingService: emailSendingService,
	}
}
