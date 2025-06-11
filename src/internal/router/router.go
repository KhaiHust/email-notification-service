package router

import (
	"github.com/KhaiHust/email-notification-service/internal/controllers"
	"github.com/KhaiHust/email-notification-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/web/actuator"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App                    *golib.App
	Engine                 *gin.Engine
	Actuator               *actuator.Endpoint
	APIKeyMiddleware       *middleware.APIKeyMiddleware
	EmailSendingController *controllers.EmailSendingController
}

func RegisterGinRouters(p RegisterRoutersIn) {
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))
	group.Use(p.APIKeyMiddleware.AuthenticationMiddlewareHandle())
	group.POST("/v1/email-request/send", p.EmailSendingController.SendEmailRequest)
}
