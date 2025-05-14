package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailTrackingController struct {
	emailTrackingService service.IEmailTrackingService
}

func (e *EmailTrackingController) OpenEmailTracking(c *gin.Context) {
	trackingID := c.Query(constant.QueryParamTrackingID)
	if trackingID == "" {
		log.Error(c, "[EmailTrackingController] trackingID is empty")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return

	}
	err := e.emailTrackingService.OpenEmailTracking(c, trackingID)
	if err != nil {
		log.Error(c, "[EmailTrackingController] Error when open email tracking", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// Write 1x1 transparent PNG
	c.Writer.Write([]byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0x0F, 0x00, 0x01,
		0x05, 0x01, 0x02, 0xA7, 0x69, 0x33, 0x00, 0x00,
		0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42,
		0x60, 0x82,
	})
}
func NewEmailTrackingController(
	emailTrackingService service.IEmailTrackingService,
) *EmailTrackingController {
	return &EmailTrackingController{
		emailTrackingService: emailTrackingService,
	}
}
