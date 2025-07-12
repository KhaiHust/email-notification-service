package bootstrap

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/client"
	strategyAdapterImpl "github.com/KhaiHust/email-notification-service/adapter/http/strategy/impl"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/adapter/publisher"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres"
	"github.com/KhaiHust/email-notification-service/adapter/service/thirdparty"
	"github.com/KhaiHust/email-notification-service/core/helper"
	middlewareCore "github.com/KhaiHust/email-notification-service/core/middleware"
	"github.com/KhaiHust/email-notification-service/core/middleware/web"
	"github.com/KhaiHust/email-notification-service/core/msg"
	coreProperties "github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/controller"
	"github.com/KhaiHust/email-notification-service/public/middleware"
	"github.com/KhaiHust/email-notification-service/public/router"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/golibs-starter/golib"
	golibdata "github.com/golibs-starter/golib-data"
	"github.com/golibs-starter/golib-data/datasource"
	golibgin "github.com/golibs-starter/golib-gin"
	golibsec "github.com/golibs-starter/golib-security"
	"github.com/golibs-starter/golib/log"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rafaelhl/gorm-newrelic-telemetry-plugin/telemetry"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net/http"
)

func All() fx.Option {
	return fx.Options(
		golib.AppOpt(),
		golib.PropertiesOpt(),
		golib.LoggingOpt(),
		golib.EventOpt(),
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		golib.ActuatorEndpointOpt(),
		golib.HttpRequestLogOpt(),
		web.CORSOpt(),
		middlewareCore.NewRelicOpt(),
		// Http security auto properties and authentication filters
		//golibsec.HttpSecurityOpt(),
		fx.Invoke(func(httpClient *http.Client) *http.Client {
			httpClient.Transport = newrelic.NewRoundTripper(http.DefaultTransport)
			return httpClient
		},
		),
		golib.HttpClientOpt(),
		golibsec.SecuredHttpClientOpt(),
		golibsec.HttpSecurityOpt(),
		golibsec.JwtAuthFilterOpt(),

		// Provide datasource auto properties
		golibdata.RedisOpt(),
		golibdata.DatasourceOpt(),
		fx.Invoke(func(conn *gorm.DB, properties *datasource.Properties) {
			err := conn.Use(telemetry.NewNrTracer(properties.Database,
				properties.Host, properties.Driver))
			if err != nil {
				log.Error("Failed to initialize New Relic telemetry plugin", err)
			}
		}),

		msg.KafkaCommonOpt(),
		msg.KafkaAdminOpt(),
		msg.KafkaProducerOpt(),

		// Provide all application properties
		golib.ProvideProps(properties.NewGmailProviderProperties),
		golib.ProvideProps(properties.NewOutlookProviderProperties),
		golib.ProvideProps(coreProperties.NewAuthProperties),
		golib.ProvideProps(coreProperties.NewBatchProperties),
		golib.ProvideProps(coreProperties.NewEncryptProperties),
		golib.ProvideProps(coreProperties.NewTrackingProperties),
		// Provide port's implements
		fx.Provide(
			fx.Annotate(client.NewGmailProviderAdapter, fx.ResultTags(`group:"emailProviderImpl"`)),
			fx.Annotate(client.NewOutlookProviderAdapter, fx.ResultTags(`group:"emailProviderImpl"`)),
			strategyAdapterImpl.NewEmailProviderAdapter,
		),

		fx.Provide(postgres.NewWorkspaceRepositoryAdapter),
		fx.Provide(postgres.NewEmailProviderRepositoryAdapter),
		fx.Provide(postgres.NewDatabaseTransactionAdapter),
		fx.Provide(postgres.NewUserRepositoryAdapter),
		fx.Provide(postgres.NewWorkspaceUserRepositoryAdapter),
		fx.Provide(postgres.NewEmailTemplateRepositoryAdapter),
		fx.Provide(postgres.NewEmailRequestRepositoryAdapter),
		fx.Provide(postgres.NewApiKeyRepositoryAdapter),
		fx.Provide(postgres.NewEmailLogRepositoryAdapter),
		fx.Provide(postgres.NewWebhookRepositoryAdapter),

		fx.Provide(publisher.NewEventPublisherAdapter),

		// Provide third-party services
		fx.Provide(thirdparty.NewRedisService),
		fx.Provide(client.NewWebhookServiceAdapter),

		// Provide use cases
		fx.Provide(usecase.NewDatabaseTransactionUseCase),
		fx.Provide(usecase.NewGetEmailProviderUseCase),
		fx.Provide(usecase.NewCreateEmailProviderUseCase),
		fx.Provide(usecase.NewGetWorkspaceUseCase),
		fx.Provide(usecase.NewCreateUserUseCase),
		fx.Provide(usecase.NewGetUserUseCase),
		fx.Provide(usecase.NewHashPasswordUseCase),
		fx.Provide(usecase.NewLoginUsecase),
		fx.Provide(usecase.NewCreateWorkspaceUseCase),
		fx.Provide(usecase.NewCreateApiKeyUseCase),
		fx.Provide(usecase.NewEncryptUseCase),
		fx.Provide(usecase.NewValidateAccessWorkspaceUsecase),
		fx.Provide(usecase.NewCreateTemplateUseCase),
		fx.Provide(usecase.NewCreateEmailRequestUsecase),
		fx.Provide(usecase.NewEmailSendingUsecase),
		fx.Provide(usecase.NewUpdateEmailRequestUsecase),
		fx.Provide(usecase.NewUpdateEmailProviderUseCase),
		fx.Provide(usecase.NewGetEmailTemplateUseCase),
		fx.Provide(usecase.NewGetEmailRequestUsecase),
		fx.Provide(usecase.NewGetApiKeyUseCase),
		fx.Provide(usecase.NewUpdateEmailTemplateUseCase),
		fx.Provide(usecase.NewEmailTrackingUsecase),
		fx.Provide(usecase.NewGetEmailLogUsecase),
		fx.Provide(usecase.NewAnalyticUsecase),
		fx.Provide(usecase.NewDeleteTemplateUseCase),
		fx.Provide(usecase.NewCreateWebhookUseCase),
		fx.Provide(usecase.NewWebhookUsecase),

		// Provide services
		fx.Provide(service.NewEmailProviderService),
		fx.Provide(service.NewUserService),
		fx.Provide(service.NewWorkspaceService),
		fx.Provide(service.NewEmailTemplateService),
		fx.Provide(service.NewEmailSendingService),
		fx.Provide(service.NewApiKeyService),
		fx.Provide(service.NewEmailRequestService),
		fx.Provide(service.NewEmailTrackingService),
		fx.Provide(service.NewEmailLogService),
		fx.Provide(service.NewAnalyticService),
		fx.Provide(service.NewWebhookService),

		//Provide controllers
		fx.Provide(helper.NewCustomValidate),
		fx.Provide(controller.NewBaseController),
		fx.Provide(controller.NewEmailProviderController),
		fx.Provide(controller.NewUserController),
		fx.Provide(controller.NewWorkspaceController),
		fx.Provide(controller.NewEmailTemplateController),
		fx.Provide(controller.NewEmailSendingController),
		fx.Provide(controller.NewApiKeyController),
		fx.Provide(controller.NewEmailRequestController),
		fx.Provide(controller.NewEmailTrackingController),
		fx.Provide(controller.NewEmailLogController),
		fx.Provide(controller.NewAnalyticController),
		fx.Provide(controller.NewWebhookController),

		golibgin.GinHttpServerOpt(),
		fx.Provide(middleware.NewWorkspaceAccessMiddleware),
		fx.Invoke(router.RegisterGinRouters),

		//Graceful shutdown
		golibgin.OnStopHttpServerOpt(),
	)
}
