package bootstrap

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/client"
	strategyAdapterImpl "github.com/KhaiHust/email-notification-service/adapter/http/strategy/impl"
	adapterProperties "github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/adapter/publisher"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres"
	"github.com/KhaiHust/email-notification-service/adapter/service/thirdparty"
	"github.com/KhaiHust/email-notification-service/core/helper"
	middlewareCore "github.com/KhaiHust/email-notification-service/core/middleware"
	"github.com/KhaiHust/email-notification-service/core/msg"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/internal/controllers"
	"github.com/KhaiHust/email-notification-service/internal/middleware"
	"github.com/KhaiHust/email-notification-service/internal/router"
	"github.com/KhaiHust/email-notification-service/internal/services"
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
		golib.HttpClientOpt(),
		fx.Invoke(func(httpClient *http.Client) *http.Client {
			httpClient.Transport = newrelic.NewRoundTripper(http.DefaultTransport)
			return httpClient
		},
		),
		golibsec.SecuredHttpClientOpt(),
		middlewareCore.NewRelicOpt(),
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

		//provide properties
		golib.ProvideProps(properties.NewBatchProperties),
		golib.ProvideProps(properties.NewEncryptProperties),
		golib.ProvideProps(properties.NewTrackingProperties),
		golib.ProvideProps(adapterProperties.NewGmailProviderProperties),
		//provider adapter
		fx.Provide(postgres.NewEmailLogRepositoryAdapter),
		fx.Provide(postgres.NewEmailRequestRepositoryAdapter),
		fx.Provide(postgres.NewEmailProviderRepositoryAdapter),
		fx.Provide(postgres.NewEmailTemplateRepositoryAdapter),
		fx.Provide(postgres.NewWorkspaceRepositoryAdapter),
		fx.Provide(postgres.NewDatabaseTransactionAdapter),
		fx.Provide(postgres.NewApiKeyRepositoryAdapter),

		//provider strategy
		fx.Provide(strategyAdapterImpl.NewEmailProviderAdapter),
		fx.Provide(client.NewGmailProviderAdapter),

		//provider service
		fx.Provide(publisher.NewEventPublisherAdapter),
		fx.Provide(thirdparty.NewRedisService),

		//provider usecase
		fx.Provide(usecase.NewGetEmailProviderUseCase),
		fx.Provide(usecase.NewGetEmailTemplateUseCase),
		fx.Provide(usecase.NewUpdateEmailProviderUseCase),
		fx.Provide(usecase.NewCreateEmailRequestUsecase),
		fx.Provide(usecase.NewDatabaseTransactionUseCase),
		fx.Provide(usecase.NewEncryptUseCase),
		fx.Provide(usecase.NewGetEmailRequestUsecase),
		fx.Provide(usecase.NewGetWorkspaceUseCase),
		fx.Provide(usecase.NewGetApiKeyUseCase),
		fx.Provide(usecase.NewValidateApiKeyUsecase),
		fx.Provide(usecase.NewUpdateEmailRequestUsecase),

		fx.Provide(usecase.NewEmailSendingUsecase),
		//provider services
		fx.Provide(services.NewEmailSendingService),

		//provider controllers
		fx.Provide(helper.NewCustomValidate),
		fx.Provide(controllers.NewBaseController),
		fx.Provide(controllers.NewEmailSendingController),

		golibgin.GinHttpServerOpt(),
		fx.Provide(middleware.NewAPIKeyMiddleware),
		fx.Invoke(router.RegisterGinRouters),

		golibgin.OnStopHttpServerOpt(),
	)

}
