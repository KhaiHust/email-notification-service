package controllers

import (
	"github.com/KhaiHust/email-notification-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailSendingController struct {
	base                *BaseController
	emailSendingService services.IEmailSendingService
}

func (esc *EmailSendingController) SendEmailRequest(c *gin.Context) {
	log.Info(c, "Received email sending request")
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
