package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
)

const (
	SubjectKey = "subject"
	BodyKey    = "body"
)

type IEmailSendingUsecase interface {
	SendBatches(ctx context.Context, providerID int64, emailRequests []*entity.EmailRequestEntity, req *request.EmailSendingRequestDto) error
	ProcessSendingEmails(ctx context.Context, workspaceID int64, req *request.EmailSendingRequestDto) error
	SendEmailByTask(ctx context.Context, emailRequest *entity.EmailRequestEntity) error
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
	emailLogRepositoryPort     port.IEmailLogRepositoryPort
}

func (e EmailSendingUsecase) SendEmailByTask(ctx context.Context, emailRequest *entity.EmailRequestEntity) error {

	decryptPayload, err := e.encryptUseCase.DecryptDataEmailRequest(ctx, string(emailRequest.Data))
	if err != nil {
		log.Error(ctx, "Error when decrypt email request data", err)
		return err
	}
	payloadMap := make(map[string]map[string]string)
	if err := json.Unmarshal([]byte(decryptPayload), &payloadMap); err != nil {
		log.Error(ctx, "Error when unmarshal email request data", err)
		return err
	}
	template, err := e.getEmailTemplateUseCase.GetTemplateByID(ctx, emailRequest.TemplateId)
	if err != nil {
		log.Error(ctx, "Error when get email template by id", err)
		return err
	}
	emailProvider, err := e.getEmailProviderUseCase.GetEmailProviderByWorkspaceID(ctx, emailRequest.WorkspaceID)
	if err != nil {
		log.Error(ctx, "Error when get email provider by id", err)
		return err
	}
	var trackingID string
	if emailRequest.TrackingID != "" {
		trackingID, err = e.encryptUseCase.EncryptTrackingID(ctx, emailRequest.TrackingID)
		if err != nil {
			return err
		}
	}
	dataSending := &request.EmailDataDto{
		EmailRequestID: emailRequest.ID,
		Subject:        utils.FillTemplate(template.Subject, payloadMap[SubjectKey]),
		Body: fmt.Sprintf(`<html><body>%s<br><img src="%s" width="100" height="100"  /></body></html>`,
			utils.FillTemplate(template.Body, payloadMap[BodyKey]),
			utils.GenerateTrackingURL(e.trackingProperties.BaseUrl, trackingID),
		),
		To: emailRequest.Recipient,
	}
	var sendErr error
	if err := e.emailProviderPort.Send(ctx, emailProvider, dataSending); err != nil {
		log.Error(ctx, fmt.Sprintf("Error when send email to %s", dataSending.To), err)
		if errors.Is(err, common.ErrUnauthorized) {
			log.Warn(ctx, fmt.Sprintf("401 detected for %s. Refreshing token...", dataSending.To))
			// Try to refresh the email provider's OAuth info
			if refreshed, refreshErr := e.updateEmailProviderUseCase.UpdateOAuthInfoByRefreshToken(ctx, emailProvider); refreshErr != nil {
				log.Error(ctx, "Error when refresh email provider OAuth info", refreshErr)
				sendErr = refreshErr
			} else {
				emailProvider = refreshed
				// Retry sending the email after refreshing the token
				if retryErr := e.emailProviderPort.Send(ctx, emailProvider, dataSending); retryErr != nil {
					log.Error(ctx, fmt.Sprintf("Retry failed to send email to %s", dataSending.To), retryErr)
					sendErr = retryErr
				}
			}
		} else {
			sendErr = err
		}
	}
	if sendErr != nil {
		emailRequest.Status = constant.EmailSendingStatusFailed
		emailRequest.ErrorMessage = sendErr.Error()
	} else {
		emailRequest.Status = constant.EmailSendingStatusSent
		now := time.Now()
		emailRequest.SentAt = utils.ToUnixTimeToPointer(&now)
	}
	go func() {
		ev := event.NewEventEmailRequestSync(ctx, emailRequest)
		e.eventPublisher.Publish(ev)
	}()
	return sendErr
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
		encryptPayload, err := e.buildEncryptVariableData(ctx, data)
		if err != nil {
			log.Error(ctx, "Error when build encrypt variable data", err)
			continue
		}
		emailRequestEntities = append(emailRequestEntities, &entity.EmailRequestEntity{
			RequestID:       requestID,
			TemplateId:      template.ID,
			Recipient:       data.To,
			Data:            encryptPayload,
			Status:          constant.EmailSendingStatusQueued,
			CorrelationID:   fmt.Sprintf("%d_%s", idx, data.To),
			EmailProviderID: emailProvider.ID,
			WorkspaceID:     workspaceID,
			TrackingID:      uuid.NewString(),
			SendAt:          data.SendAt,
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
	if err = e.SaveEmailLogsByBatches(ctx, tx, emailRequestEntities); err != nil {
		return err
	}
	if err = e.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Error when commit transaction", err)
		return err
	}
	commit = true
	ev := event.NewEventRequestSendingEmail(ctx, emailRequestEntities, req)
	go func() {
		if err := e.eventPublisher.SyncPublish(ctx, ev); err != nil {
			log.Error(ctx, "Error when publish event", err)
		}
	}()
	return nil
}

func (e EmailSendingUsecase) SaveEmailLogsByBatches(ctx context.Context, tx *gorm.DB, emailRequestEntities []*entity.EmailRequestEntity) error {
	emailLogsEntities := make([]*entity.EmailLogsEntity, 0, len(emailRequestEntities))
	for _, emailRequest := range emailRequestEntities {
		emailLogsEntities = append(emailLogsEntities, &entity.EmailLogsEntity{
			EmailRequestID:  emailRequest.ID,
			TemplateId:      emailRequest.TemplateId,
			Recipient:       emailRequest.Recipient,
			Status:          constant.EmailSendingStatusQueued,
			RequestID:       emailRequest.RequestID,
			WorkspaceID:     emailRequest.WorkspaceID,
			EmailProviderID: emailRequest.EmailProviderID,
			LoggedAt:        time.Now().Unix(),
		})
	}
	if _, err := e.emailLogRepositoryPort.SaveEmailLogsByBatches(ctx, tx, emailLogsEntities); err != nil {
		log.Error(ctx, "Error when save email logs", err)
		return err
	}
	return nil
}

func (e EmailSendingUsecase) SendBatches(ctx context.Context, providerID int64, emailRequests []*entity.EmailRequestEntity, req *request.EmailSendingRequestDto) error {
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

	// Prepare emailRequest to send
	dataSendings := make([]*request.EmailDataDto, 0, len(req.Datas))
	for _, emailRequest := range emailRequests {
		var trackingID string
		if emailRequest.TrackingID != "" {
			trackingID, err = e.encryptUseCase.EncryptTrackingID(ctx, emailRequest.TrackingID)
			if err != nil {
				continue
			}
		}
		decryptPayload, err := e.encryptUseCase.DecryptDataEmailRequest(ctx, string(emailRequest.Data))
		if err != nil {
			log.Error(ctx, "Error when decrypt email request data", err)
			continue
		}
		payloadMap := make(map[string]map[string]string)
		if err := json.Unmarshal([]byte(decryptPayload), &payloadMap); err != nil {
			log.Error(ctx, "Error when unmarshal email request data", err)
			continue
		}
		dataSendings = append(dataSendings, &request.EmailDataDto{
			EmailRequestID: emailRequest.ID,
			Subject:        utils.FillTemplate(template.Subject, payloadMap[SubjectKey]),
			Body: fmt.Sprintf(`<html><body>%s<br><img src="%s" width="100" height="100"  /></body></html>`,
				utils.FillTemplate(template.Body, payloadMap[BodyKey]),
				utils.GenerateTrackingURL(e.trackingProperties.BaseUrl, trackingID),
			),
			To: emailRequest.Recipient,
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
					log.Warn(ctx, fmt.Sprintf("401 detected for %v. Refreshing token...", data.To))

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
					log.Error(ctx, fmt.Sprintf("Failed to send email to %v", data.To), sendErr)
				}
				//fire event to sync email request
				emailRequest := &entity.EmailRequestEntity{
					BaseEntity: entity.BaseEntity{
						ID: data.EmailRequestID,
					},
					Recipient:    data.To,
					Status:       status,
					ErrorMessage: errMessage,
					SentAt:       sentAt,
				}
				ev := event.NewEventEmailRequestSync(ctx, emailRequest)
				e.eventPublisher.Publish(ev)

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
	emailLogRepositoryPort port.IEmailLogRepositoryPort,
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
		emailLogRepositoryPort:     emailLogRepositoryPort,
	}
}
func (e EmailSendingUsecase) buildEncryptVariableData(ctx context.Context, rawData *request.EmailSendingData) (string, error) {
	var data = make(map[string]interface{})
	data[SubjectKey] = rawData.Subject
	data[BodyKey] = rawData.Body

	payload, err := json.Marshal(data)
	if err != nil {
		log.Error(ctx, "Error when marshal payload data", err)
		return "", err
	}
	// Encrypt the payload
	encryptedPayload, err := e.encryptUseCase.EncryptDataEmailRequest(ctx, string(payload))
	if err != nil {
		log.Error(ctx, "Error when encrypt payload data", err)
		return "", err
	}
	return encryptedPayload, nil
}
