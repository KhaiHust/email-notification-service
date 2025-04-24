package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
)

type IEmailSendingUsecase interface {
	SendBatches(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error
}
type EmailSendingUsecase struct {
	getEmailProviderUseCase IGetEmailProviderUseCase
	getEmailTemplateUseCase IGetEmailTemplateUseCase
	emailProviderPort       port.IEmailProviderPort
}

func (e EmailSendingUsecase) SendBatches(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error {
	emailProvider, err := e.getEmailProviderUseCase.GetEmailProviderByID(ctx, providerID)
	if err != nil {
		log.Error(ctx, "Error when get email provider by id", err)
		return err
	}

	template, err := e.getEmailTemplateUseCase.GetTemplateByID(ctx, req.TemplateId)
	if err != nil {
		log.Error(ctx, "Error when get email template by id", err)
		return err
	}
	dataSendings := make([]*request.EmailDataDto, 0)
	for _, data := range req.Data {
		dataSend := &request.EmailDataDto{
			Subject: utils.FillTemplate(template.Subject, data.Subject),
			Body:    utils.FillTemplate(template.Body, data.Body),
			Tos:     []string{data.To},
		}
		dataSendings = append(dataSendings, dataSend)
	}

	for _, data := range dataSendings {
		err = e.emailProviderPort.Send(ctx, emailProvider, data)
		if err != nil {
			log.Error(ctx, "Error when send email", err)
			return err
		}
	}

	return nil
}

func NewEmailSendingUsecase(
	getEmailProviderUseCase IGetEmailProviderUseCase,
	getEmailTemplateUseCase IGetEmailTemplateUseCase,
	emailProviderPort port.IEmailProviderPort,
) IEmailSendingUsecase {
	return &EmailSendingUsecase{
		getEmailProviderUseCase: getEmailProviderUseCase,
		getEmailTemplateUseCase: getEmailTemplateUseCase,
		emailProviderPort:       emailProviderPort,
	}
}
