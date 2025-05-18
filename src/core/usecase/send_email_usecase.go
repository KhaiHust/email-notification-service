package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/google/uuid"
	"sync"
	"time"

	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
)

type IEmailSendingUsecase interface {
	SendBatches(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error
	ProcessSendingEmails(ctx context.Context, workspaceID int64, req *request.EmailSendingRequestDto) error
}

type EmailSendingUsecase struct {
	BatchConfig                *properties.BatchProperties
	getEmailProviderUseCase    IGetEmailProviderUseCase
	getEmailTemplateUseCase    IGetEmailTemplateUseCase
	emailProviderPort          port.IEmailProviderPort
	updateEmailProviderUseCase IUpdateEmailProviderUseCase
	eventPublisher             port.IEventPublisher
	createEmailRequestUsecase  ICreateEmailRequestUsecase
	databaseTransactionUseCase IDatabaseTransactionUseCase
	trackingProperties         *properties.TrackingProperties
	encryptUseCase             IEncryptUseCase
}

func (e EmailSendingUsecase) ProcessSendingEmails(ctx context.Context, workspaceID int64, req *request.EmailSendingRequestDto) error {
	emailProvider, err := e.getEmailProviderUseCase.GetEmailProviderByWorkspaceID(ctx, workspaceID)
	if err != nil {
		log.Error(ctx, "Error when get email provider by id", err)
		return err
	}

	template, err := e.getEmailTemplateUseCase.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		log.Error(ctx, "Error when get email template by id", err)
		return err
	}
	req.IntegrationID = emailProvider.ID
	req.TemplateID = template.ID

	emailRequestEntities := make([]*entity.EmailRequestEntity, 0, len(req.Datas))
	requestID := uuid.New().String()
	for idx, data := range req.Datas {
		emailRequestEntities = append(emailRequestEntities, &entity.EmailRequestEntity{
			RequestID:       requestID,
			TemplateId:      template.ID,
			Recipient:       data.To,
			Status:          constant.EmailSendingStatusQueued,
			CorrelationID:   fmt.Sprintf("%d_%s", idx, data.To),
			EmailProviderID: emailProvider.ID,
			WorkspaceID:     workspaceID,
			TrackingID:      uuid.NewString(),
		})
	}
	tx := e.databaseTransactionUseCase.StartTx()
	commit := false
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if !commit || err != nil {
			if errRollback := e.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Error when rollback transaction", errRollback)
			} else {
				log.Info(ctx, "Rollback transaction successfully")
			}
		}
	}()
	emailRequestEntities, err = e.createEmailRequestUsecase.CreateEmailRequestsWithTx(ctx, tx, emailRequestEntities)
	if err != nil {
		log.Error(ctx, "Error when save email request", err)
		return err
	}
	ev := event.NewEventRequestSendingEmail(ctx, emailRequestEntities, req)
	err = e.eventPublisher.SyncPublish(ctx, ev)
	if err != nil {
		log.Error(ctx, "Error when publish event", err)
		return err
	}
	if err = e.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Error when commit transaction", err)
		return err
	}
	commit = true
	return nil
}

func (e EmailSendingUsecase) SendBatches(ctx context.Context, providerID int64, req *request.EmailSendingRequestDto) error {
	emailProvider, err := e.getEmailProviderUseCase.GetEmailProviderByID(ctx, providerID)
	if err != nil {
		log.Error(ctx, "Error when get email provider by id", err)
		return err
	}

	template, err := e.getEmailTemplateUseCase.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		log.Error(ctx, "Error when get email template by id", err)
		return err
	}

	// Prepare data to send
	dataSendings := make([]*request.EmailDataDto, 0, len(req.Datas))
	for _, data := range req.Datas {
		var trackingID string
		if data.TrackingID != "" {
			trackingID, err = e.encryptUseCase.EncryptTrackingID(ctx, data.TrackingID)
			if err != nil {
				continue
			}
		}
		dataSendings = append(dataSendings, &request.EmailDataDto{
			EmailRequestID: data.EmailRequestID,
			Subject:        utils.FillTemplate(template.Subject, data.Subject),
			Body: fmt.Sprintf(`<html><body>%s<br><img src="%s" width="100" height="100"  /></body></html>`,
				utils.FillTemplate(template.Body, data.Body),
				utils.GenerateTrackingURL(e.trackingProperties.BaseUrl, trackingID),
			),
			Tos: []string{data.To},
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
				sendErr := e.emailProviderPort.Send(ctx, emailProvider, data)
				status := constant.EmailSendingStatusSent
				var errMessage string
				timeNow := time.Now()
				sentAt := utils.ToUnixTimeToPointer(&timeNow)

				if errors.Is(sendErr, common.ErrUnauthorized) {
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
						sendErr = e.emailProviderPort.Send(ctx, emailProvider, data)
					}
				}

				if sendErr != nil {
					status = constant.EmailSendingStatusFailed
					errMessage = sendErr.Error()
					log.Error(ctx, fmt.Sprintf("Failed to send email to %v", data.Tos), sendErr)
				}
				//fire event to sync email request
				for _, to := range data.Tos {
					emailRequest := &entity.EmailRequestEntity{
						BaseEntity: entity.BaseEntity{
							ID: data.EmailRequestID,
						},
						Recipient:    to,
						Status:       status,
						ErrorMessage: errMessage,
						SentAt:       sentAt,
					}
					ev := event.NewEventEmailRequestSync(ctx, emailRequest)
					e.eventPublisher.Publish(ev)
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
	eventPublisher port.IEventPublisher,
	createEmailRequestUsecase ICreateEmailRequestUsecase,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	trackingProperties *properties.TrackingProperties,
	encryptUseCase IEncryptUseCase,
) IEmailSendingUsecase {
	return &EmailSendingUsecase{
		BatchConfig:                batchConfig,
		getEmailProviderUseCase:    getEmailProviderUseCase,
		getEmailTemplateUseCase:    getEmailTemplateUseCase,
		emailProviderPort:          emailProviderPort,
		updateEmailProviderUseCase: updateEmailProviderUseCase,
		eventPublisher:             eventPublisher,
		createEmailRequestUsecase:  createEmailRequestUsecase,
		databaseTransactionUseCase: databaseTransactionUseCase,
		trackingProperties:         trackingProperties,
		encryptUseCase:             encryptUseCase,
	}
}
