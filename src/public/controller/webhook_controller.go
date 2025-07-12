package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
	"strconv"
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
func (w *WebhookController) GetDetail(c *gin.Context) {
	workspaceID := w.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is required")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	webhookID, err := strconv.ParseInt(c.Param(constant.ParamWebhookId), 10, 64)
	if err != nil {
		log.Error(c, "Failed to get webhook ID from context", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	webhookResponse, err := w.GetWebhookDetail(c, workspaceID, webhookID)
	if err != nil {
		log.Error(c, "Failed to get webhook detail", "error", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, webhookResponse)
}
func (w *WebhookController) GetAll(c *gin.Context) {
	workspaceID := w.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is required")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	webhooks, err := w.GetAllWebhooksByWorkspaceID(c, workspaceID)
	if err != nil {
		log.Error(c, "Failed to get all webhooks", "error", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, webhooks)
}
func (w *WebhookController) Delete(c *gin.Context) {
	workspaceID := w.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is required")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	webhookID, err := strconv.ParseInt(c.Param(constant.ParamWebhookId), 10, 64)
	if err != nil {
		log.Error(c, "Failed to get webhook ID from context", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err = w.DeleteWebhook(c, workspaceID, webhookID); err != nil {
		log.Error(c, "Failed to delete webhook", "error", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, nil)
}
func (w *WebhookController) Update(c *gin.Context) {
	workspaceID := w.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is required")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	webhookID, err := strconv.ParseInt(c.Param(constant.ParamWebhookId), 10, 64)
	if err != nil {
		log.Error(c, "Failed to get webhook ID from context", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	var req request.UpdateWebhookRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "Failed to bind request", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err = w.validator.Struct(&req); err != nil {
		log.Error(c, "Validation failed", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	webhookResponse, err := w.UpdateWebhook(c, workspaceID, webhookID, &req)
	if err != nil {
		log.Error(c, "Failed to update webhook", "error", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, webhookResponse)
}
func (w *WebhookController) Test(c *gin.Context) {
	workspaceID := w.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Workspace ID is required")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
	}
	webhookID, err := strconv.ParseInt(c.Param(constant.ParamWebhookId), 10, 64)
	if err != nil {
		log.Error(c, "Failed to get webhook ID from context", "error", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err = w.TestWebhook(c, workspaceID, webhookID); err != nil {
		log.Error(c, "Failed to test webhook", "error", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, nil)
}
func NewWebhookController(base *BaseController, webhookService service.IWebhookService) *WebhookController {
	return &WebhookController{
		BaseController:  base,
		IWebhookService: webhookService,
	}
}
