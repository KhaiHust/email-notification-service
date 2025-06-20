package bootstrap

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/client"
	strategyAdapterImpl "github.com/KhaiHust/email-notification-service/adapter/http/strategy/impl"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/adapter/publisher"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres"
	"github.com/KhaiHust/email-notification-service/adapter/service/thirdparty"
	"github.com/KhaiHust/email-notification-service/core/msg"
	coreProperties "github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/worker/cronjob"
	"github.com/KhaiHust/email-notification-service/worker/handler"
	"github.com/KhaiHust/email-notification-service/worker/router"
	"github.com/golibs-starter/golib"
	golibcron "github.com/golibs-starter/golib-cron"
	golibdata "github.com/golibs-starter/golib-data"
	golibgin "github.com/golibs-starter/golib-gin"
	golibmsg "github.com/golibs-starter/golib-message-bus"
	golibsec "github.com/golibs-starter/golib-security"
	"go.uber.org/fx"
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
		golibsec.SecuredHttpClientOpt(),

		msg.KafkaCommonOpt(),
		msg.KafkaAdminOpt(),
		msg.KafkaProducerOpt(),
		msg.KafkaConsumerOpt(),

		// Provide cronjob
		golibcron.Opt(),
		// Provide datasource auto properties
		golibdata.RedisOpt(),
		golibdata.DatasourceOpt(),

		// Provide all application properties
		golib.ProvideProps(properties.NewGmailProviderProperties),
		golib.ProvideProps(properties.NewOutlookProviderProperties),
		golib.ProvideProps(properties.NewGoogleCloudTaskProperties),
		golib.ProvideProps(coreProperties.NewAuthProperties),
		golib.ProvideProps(coreProperties.NewBatchProperties),
		golib.ProvideProps(coreProperties.NewEncryptProperties),
		golib.ProvideProps(coreProperties.NewTrackingProperties),
		golib.ProvideProps(coreProperties.NewTaskProperties),
		golib.ProvideProps(coreProperties.NewAppProperties),
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
		fx.Provide(postgres.NewEmailLogRepositoryAdapter),
		fx.Provide(postgres.NewWebhookRepositoryAdapter),

		fx.Provide(publisher.NewEventPublisherAdapter),
		fx.Provide(thirdparty.NewCloudTaskServiceAdapter),

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
		fx.Provide(usecase.NewValidateAccessWorkspaceUsecase),
		fx.Provide(usecase.NewCreateTemplateUseCase),
		fx.Provide(usecase.NewEmailSendingUsecase),
		fx.Provide(usecase.NewGetEmailTemplateUseCase),
		fx.Provide(usecase.NewUpdateEmailProviderUseCase),
		fx.Provide(usecase.NewEventHandlerUsecase),
		fx.Provide(usecase.NewCreateEmailRequestUsecase),
		fx.Provide(usecase.NewUpdateEmailRequestUsecase),
		fx.Provide(usecase.NewGetEmailRequestUsecase),
		fx.Provide(usecase.NewEncryptUseCase),
		fx.Provide(usecase.NewScheduleEmailUsecase),
		fx.Provide(usecase.NewEmailSendRetryUsecase),
		fx.Provide(usecase.NewWebhookUsecase),

		//provider handler
		golibmsg.ProvideConsumer(handler.NewEmailSendingRequestHandler),
		golibmsg.ProvideConsumer(handler.NewEmailRequestSyncHandler),

		//provide cronjob
		golibcron.ProvideJob(cronjob.NewEmailSendRetryCronJob),
		// Provide gin engine, register core handlers,
		// actuator endpoints and application routers
		golibgin.GinHttpServerOpt(),
		fx.Invoke(router.RegisterGinRouters),
		// Graceful shutdown.
		// OnStop hooks will run in reverse order.
		//golibgin.OnStopHttpServerOpt(),
		//golibmsg.OnStopProducerOpt(),
		golibmsg.OnStopConsumerOpt(),
	)
}
