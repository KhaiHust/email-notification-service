package bootstrap

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/client"
	strategyAdapterImpl "github.com/KhaiHust/email-notification-service/adapter/http/strategy/impl"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/adapter/publisher"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres"
	"github.com/KhaiHust/email-notification-service/adapter/service/thirdparty"
	"github.com/KhaiHust/email-notification-service/core/helper"
	"github.com/KhaiHust/email-notification-service/core/middleware/web"
	coreProperties "github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/controller"
	"github.com/KhaiHust/email-notification-service/public/middleware"
	"github.com/KhaiHust/email-notification-service/public/router"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/golibs-starter/golib"
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
		web.CORSOpt(),
		// Http security auto config and authentication filters
		//golibsec.HttpSecurityOpt(),

		golib.HttpClientOpt(),
		golibsec.SecuredHttpClientOpt(),
		golibsec.HttpSecurityOpt(),
		golibsec.JwtAuthFilterOpt(),

		// Provide datasource auto config
		golibdata.RedisOpt(),
		golibdata.DatasourceOpt(),
		golibmsg.KafkaCommonOpt(),
		golibmsg.KafkaAdminOpt(),
		golibmsg.KafkaProducerOpt(),

		// Provide all application properties
		golib.ProvideProps(properties.NewGmailProviderProperties),
		golib.ProvideProps(coreProperties.NewAuthProperties),
		golib.ProvideProps(coreProperties.NewBatchProperties),
		golib.ProvideProps(coreProperties.NewEncryptProperties),
		// Provide port's implements
		fx.Provide(client.NewGmailProviderAdapter),
		fx.Provide(strategyAdapterImpl.NewEmailProviderAdapter),
		fx.Provide(postgres.NewWorkspaceRepositoryAdapter),
		fx.Provide(postgres.NewEmailProviderRepositoryAdapter),
		fx.Provide(postgres.NewDatabaseTransactionAdapter),
		fx.Provide(postgres.NewUserRepositoryAdapter),
		fx.Provide(postgres.NewWorkspaceUserRepositoryAdapter),
		fx.Provide(postgres.NewEmailTemplateRepositoryAdapter),
		fx.Provide(postgres.NewEmailRequestRepositoryAdapter),
		fx.Provide(postgres.NewApiKeyRepositoryAdapter),

		fx.Provide(publisher.NewEventPublisherAdapter),

		// Provide third-party services
		fx.Provide(thirdparty.NewRedisService),

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
		// Provide services
		fx.Provide(service.NewEmailProviderService),
		fx.Provide(service.NewUserService),
		fx.Provide(service.NewWorkspaceService),
		fx.Provide(service.NewEmailTemplateService),
		fx.Provide(service.NewEmailSendingService),
		fx.Provide(service.NewApiKeyService),
		fx.Provide(service.NewEmailRequestService),

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

		golibgin.GinHttpServerOpt(),
		fx.Provide(middleware.NewWorkspaceAccessMiddleware),
		fx.Invoke(router.RegisterGinRouters),

		//Graceful shutdown
		golibgin.OnStopHttpServerOpt(),
	)
}
