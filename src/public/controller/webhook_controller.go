package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type WebhookController struct {
	*BaseController
	service.IWebhookService
}

func (w *WebhookController) Create(c *gin.Context) {
	workspaceID := w.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is required")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	var req request.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "Failed to bind request", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := w.validator.Struct(&req); err != nil {
		log.Error(c, "Validation failed", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	webhookResponse, err := w.CreateWebhook(c, workspaceID, &req)
	if err != nil {
		log.Error(c, "Failed to create webhook", "error", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, webhookResponse)
}
func NewWebhookController(base *BaseController, webhookService service.IWebhookService) *WebhookController {
	return &WebhookController{
		BaseController:  base,
		IWebhookService: webhookService,
	}
}
