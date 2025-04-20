package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailTemplateController struct {
	*BaseController
	emailTemplateService service.IEmailTemplateService
}

func (e *EmailTemplateController) CreateTemplate(c *gin.Context) {
	userId, err := e.GetUserIDFromContext(c)
	if err != nil {
		log.Error(c, "Error when get user id from context", err)
		apihelper.AbortErrorHandle(c, common.ErrForbidden)
		return
	}
	workspaceID := c.Param(constant.ParamWorkspaceCode)
	if workspaceID == "" {
		log.Error(c, "Error when get workspace id from context", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	var req request.CreateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "Error when bind json", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := e.validator.Struct(&req); err != nil {
		log.Error(c, "Error when validate request", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	emailTemplate, err := e.emailTemplateService.CreateTemplate(c, userId, workspaceID, &req)
	if err != nil {
		log.Error(c, "Error when create email template", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, response.ToEmailTemplateResponse(emailTemplate))

}
func NewEmailTemplateController(
	emailTemplateService service.IEmailTemplateService,
	base *BaseController,
) *EmailTemplateController {
	return &EmailTemplateController{
		BaseController:       base,
		emailTemplateService: emailTemplateService,
	}
}
