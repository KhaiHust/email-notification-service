package router

import (
	"github.com/KhaiHust/email-notification-service/internal/controllers"
	"github.com/KhaiHust/email-notification-service/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/web/actuator"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App                     *golib.App
	Engine                  *gin.Engine
	Actuator                *actuator.Endpoint
	APIKeyMiddleware        *middleware.APIKeyMiddleware
	EmailSendingController  *controllers.EmailSendingController
	EmailTrackingController *controllers.EmailTrackingController
	NewRelic                *newrelic.Application
}

func RegisterGinRouters(p RegisterRoutersIn) {
	p.Engine.Use(nrgin.Middleware(p.NewRelic))
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))
	v1Tracking := group.Group("/v1/tracking")
	{
		v1Tracking.GET("", p.EmailTrackingController.OpenEmailTracking)
	}
	v1Task := group.Group("/v1/tasks")
	{
		v1Task.POST("/send", p.EmailSendingController.SendEmailByTask)
	}
	group.Use(p.APIKeyMiddleware.AuthenticationMiddlewareHandle())
	v1EmailRequest := group.Group("/v1/email-requests")
	{
		v1EmailRequest.POST("/send", p.EmailSendingController.SendEmailRequest)
	}
}
