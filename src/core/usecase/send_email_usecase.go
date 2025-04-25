package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"sync"

	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
)

type IEmailSendingUsecase interface {
	SendBatches(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error
}

type EmailSendingUsecase struct {
	BatchConfig                *properties.BatchProperties
	getEmailProviderUseCase    IGetEmailProviderUseCase
	getEmailTemplateUseCase    IGetEmailTemplateUseCase
	emailProviderPort          port.IEmailProviderPort
	updateEmailProviderUseCase IUpdateEmailProviderUseCase
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

	// Prepare data to send
	dataSendings := make([]*request.EmailDataDto, 0, len(req.Data))
	for _, data := range req.Data {
		dataSendings = append(dataSendings, &request.EmailDataDto{
			Subject: utils.FillTemplate(template.Subject, data.Subject),
			Body:    utils.FillTemplate(template.Body, data.Body),
			Tos:     []string{data.To},
		})
	}

	// Setup worker pool
	numWorkers := e.BatchConfig.NumOfWorkers
	jobs := make(chan *request.EmailDataDto)
	var wg sync.WaitGroup

	var (
		refreshOnce sync.Once
		refreshErr  error
		mu          sync.Mutex
	)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range jobs {
				err := e.emailProviderPort.Send(ctx, emailProvider, data)
				if errors.Is(err, common.ErrUnauthorized) {
					log.Warn(ctx, fmt.Sprintf("401 detected for %v. Refreshing token...", data.Tos))

					refreshOnce.Do(func() {
						refreshed, err := e.updateEmailProviderUseCase.UpdateOAuthInfoByRefreshToken(ctx, emailProvider)
						if err != nil {
							refreshErr = err
							log.Error(ctx, "Token refresh failed", err)
							return
						}
						mu.Lock()
						emailProvider = refreshed
						mu.Unlock()
					})

					if refreshErr == nil {
						err = e.emailProviderPort.Send(ctx, emailProvider, data)
					}
				}

				if err != nil {
					log.Error(ctx, fmt.Sprintf("Failed to send email to %v", data.Tos), err)
				}
			}
		}()
	}

	for _, data := range dataSendings {
		jobs <- data
	}
	close(jobs)
	wg.Wait()

	return nil
}

func NewEmailSendingUsecase(
	batchConfig *properties.BatchProperties,
	getEmailProviderUseCase IGetEmailProviderUseCase,
	getEmailTemplateUseCase IGetEmailTemplateUseCase,
	emailProviderPort port.IEmailProviderPort,
	updateEmailProviderUseCase IUpdateEmailProviderUseCase,
) IEmailSendingUsecase {
	return &EmailSendingUsecase{
		BatchConfig:                batchConfig,
		getEmailProviderUseCase:    getEmailProviderUseCase,
		getEmailTemplateUseCase:    getEmailTemplateUseCase,
		emailProviderPort:          emailProviderPort,
		updateEmailProviderUseCase: updateEmailProviderUseCase,
	}
}
