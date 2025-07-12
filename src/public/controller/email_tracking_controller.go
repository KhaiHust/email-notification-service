package controller

import (
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
	"os"
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
	cwd, _ := os.Getwd()
	c.File(fmt.Sprintf("%s/resource/__files/logo.png", cwd))
}
func NewEmailTrackingController(
	emailTrackingService service.IEmailTrackingService,
) *EmailTrackingController {
	return &EmailTrackingController{
		emailTrackingService: emailTrackingService,
	}
}
