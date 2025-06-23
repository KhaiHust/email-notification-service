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
	"mime"
	"sync"
	"time"

	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
)

const (
	SubjectKey       = "subject"
	BodyKey          = "body"
	PrefixCacheLease = "email_request_lease_%d"
	LeaseDuration    = 60 // seconds
)

type IEmailSendingUsecase interface {
	SendBatches(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error
	ProcessSendingEmails(ctx context.Context, workspaceID int64, req *request.EmailSendingRequestDto) (string, error)
	SendEmailByTask(ctx context.Context, emailRequest *entity.EmailRequestEntity) error
	SendSyncs(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error
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
	updateEmailRequestUsecase  IUpdateEmailRequestUsecase
	redisPort                  port.IRedisPort
}

func (e EmailSendingUsecase) SendEmailByTask(ctx context.Context, emailRequest *entity.EmailRequestEntity) error {

	decryptPayload, err := e.encryptUseCase.DecryptDataEmailRequest(ctx, emailRequest.Data)
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
	emailProvider, err := e.getEmailProviderUseCase.GetEmailProviderByIDAndWorkspaceID(ctx, emailRequest.EmailProviderID, emailRequest.WorkspaceID)
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

func (e EmailSendingUsecase) ProcessSendingEmails(ctx context.Context, workspaceID int64, req *request.EmailSendingRequestDto) (string, error) {
	var emailProvider *entity.EmailProviderEntity
	if req.Provider != nil {
		emailProvider = req.Provider
	} else {
		emailProviderEntity, err := e.getEmailProviderUseCase.GetEmailProviderByIDAndWorkspaceID(ctx, req.ProviderID, workspaceID)
		if err != nil {
			log.Error(ctx, "Error when get email provider by id", err)
			return "", err
		}
		emailProvider = emailProviderEntity
	}
	template, err := e.getEmailTemplateUseCase.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		log.Error(ctx, "Error when get email template by id", err)
		return "", err
	}
	req.ProviderID = emailProvider.ID
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
	if len(emailRequestEntities) == 0 {
		log.Warn(ctx, "No email requests to process")
		return requestID, fmt.Errorf("request ID: %s, No email requests to process", requestID)
	}
	emailRequestEntities, err = e.createEmailRequestUsecase.CreateEmailRequestsWithTx(ctx, tx, emailRequestEntities)
	if err != nil {
		log.Error(ctx, "Error when save email request", err)
		return "", err
	}
	if err = e.SaveEmailLogsByBatches(ctx, tx, emailRequestEntities); err != nil {
		return "", err
	}
	if err = e.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Error when commit transaction", err)
		return "", err
	}
	commit = true
	ev := event.NewEventRequestSendingEmail(ctx, emailRequestEntities)
	go func() {
		if err := e.eventPublisher.SyncPublish(ctx, ev); err != nil {
			log.Error(ctx, "Error when publish event", err)
		}
	}()
	return requestID, nil
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

func (e EmailSendingUsecase) SendBatches(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error {
	emailProviderMap, templateMap, err := e.buildTemplateAndProviderMap(ctx, emailRequests)
	if err != nil {
		return err
	}
	// Prepare emailRequest to send
	dataSendings := make([]*request.EmailDataDto, 0, len(emailRequests))
	for _, emailRequest := range emailRequests {
		emailProvider, ok := emailProviderMap[emailRequest.EmailProviderID]
		if !ok {
			log.Error(ctx, fmt.Sprintf("Email provider with ID %d not found for email request %d", emailRequest.EmailProviderID, emailRequest.ID))
			continue
		}
		template, ok := templateMap[emailRequest.TemplateId]
		if !ok {
			log.Error(ctx, fmt.Sprintf("Template with ID %d not found for email request %d", emailRequest.TemplateId, emailRequest.ID))
			continue
		}
		var trackingID string
		if emailRequest.TrackingID != "" {
			trackingID, err = e.encryptUseCase.EncryptTrackingID(ctx, emailRequest.TrackingID)
			if err != nil {
				continue
			}
		}
		decryptPayload, err := e.encryptUseCase.DecryptDataEmailRequest(ctx, emailRequest.Data)
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
			Body: fmt.Sprintf(`<html><body>%s<br><img src="%s" width="1" height="1"  /></body></html>`,
				utils.FillTemplate(template.Body, payloadMap[BodyKey]),
				utils.GenerateTrackingURL(e.trackingProperties.BaseUrl, trackingID),
			),
			To:       emailRequest.Recipient,
			Provider: emailProvider,
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
				sendErr := e.emailProviderPort.Send(ctx, data.Provider, data)
				status := constant.EmailSendingStatusSent
				var errMessage string
				timeNow := time.Now()
				sentAt := utils.ToUnixTimeToPointer(&timeNow)

				if errors.Is(sendErr, common.ErrUnauthorized) {
					log.Warn(ctx, fmt.Sprintf("401 detected for %v. Refreshing token...", data.To))

					refreshOnce.Do(func() {
						refreshed, err := e.updateEmailProviderUseCase.UpdateOAuthInfoByRefreshToken(ctx, data.Provider)
						if err != nil {
							refreshErr = err
							log.Error(ctx, "Token refresh failed", err)
							return
						}
						mu.Lock()
						data.Provider = refreshed
						mu.Unlock()
					})

					if refreshErr == nil {
						sendErr = e.emailProviderPort.Send(ctx, data.Provider, data)
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
	updateEmailRequestUsecase IUpdateEmailRequestUsecase,
	redisPort port.IRedisPort,
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
		updateEmailRequestUsecase:  updateEmailRequestUsecase,
		redisPort:                  redisPort,
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
func (e EmailSendingUsecase) SendSyncs(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error {
	emailProviderMap, templateMap, err := e.buildTemplateAndProviderMap(ctx, emailRequests)
	if err != nil {
		return err
	}

	dataSendings := make([]*request.EmailDataDto, 0, len(emailRequests))
	for idx, emailRequest := range emailRequests {
		now := time.Now().Unix()
		emailRequest.SentAt = &now

		emailProvider, ok := emailProviderMap[emailRequest.EmailProviderID]
		if !ok {
			log.Error(ctx, fmt.Sprintf("Email provider with ID %d not found for email request %d", emailRequest.EmailProviderID, emailRequest.ID))
			emailRequests[idx].Status = constant.EmailSendingStatusFailed
			emailRequests[idx].ErrorMessage = fmt.Sprintf("Email provider with ID %d not found", emailRequest.EmailProviderID)
			continue
		}

		template, ok := templateMap[emailRequest.TemplateId]
		if !ok {
			log.Error(ctx, fmt.Sprintf("Template with ID %d not found for email request %d", emailRequest.TemplateId, emailRequest.ID))
			emailRequests[idx].Status = constant.EmailSendingStatusFailed
			emailRequests[idx].ErrorMessage = fmt.Sprintf("Template with ID %d not found", emailRequest.TemplateId)
			continue
		}

		var trackingID string
		if emailRequest.TrackingID != "" {
			trackingID, err = e.encryptUseCase.EncryptTrackingID(ctx, emailRequest.TrackingID)
			if err != nil {
				log.Error(ctx, "Error when encrypt tracking ID", err)
				emailRequests[idx].Status = constant.EmailSendingStatusFailed
				emailRequests[idx].ErrorMessage = "Error when encrypt tracking ID"
				continue
			}
		}

		decryptPayload, err := e.encryptUseCase.DecryptDataEmailRequest(ctx, emailRequest.Data)
		if err != nil {
			log.Error(ctx, "Error when decrypt email request data", err)
			emailRequests[idx].Status = constant.EmailSendingStatusFailed
			emailRequests[idx].ErrorMessage = "Error when decrypt email request data"
			continue
		}

		payloadMap := make(map[string]map[string]string)
		if err := json.Unmarshal([]byte(decryptPayload), &payloadMap); err != nil {
			log.Error(ctx, "Error when unmarshal email request data", err)
			emailRequests[idx].Status = constant.EmailSendingStatusFailed
			emailRequests[idx].ErrorMessage = "Error when unmarshal email request data"
			continue
		}
		subject := utils.FillTemplate(template.Subject, payloadMap[SubjectKey])
		encodedSubject := mime.BEncoding.Encode("UTF-8", subject)
		dataSendings = append(dataSendings, &request.EmailDataDto{
			EmailRequestID: emailRequest.ID,
			Subject:        encodedSubject,
			Body: fmt.Sprintf(`<html><body>%s<br><img src="%s" width="1" height="1"  /></body></html>`,
				utils.FillTemplate(template.Body, payloadMap[BodyKey]),
				utils.GenerateTrackingURL(e.trackingProperties.BaseUrl, trackingID),
			),
			To:       emailRequest.Recipient,
			Provider: emailProvider,
		})
	}

	mapEmailRequests := make(map[int64]*entity.EmailRequestEntity)
	for _, emailRequest := range emailRequests {
		mapEmailRequests[emailRequest.ID] = emailRequest
	}

	numWorkers := e.BatchConfig.NumOfWorkers
	jobs := make(chan *request.EmailDataDto, numWorkers*2)
	results := make(chan *entity.EmailRequestEntity, len(dataSendings))
	var wg sync.WaitGroup

	// Per-provider refreshOnce and refreshErr storage
	type refreshData struct {
		once *sync.Once
		err  error
		mu   sync.Mutex
	}
	providerRefreshMap := sync.Map{} // map[int64]*refreshData

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range jobs {
				// 1. Acquire lease on Redis
				leaseKey := fmt.Sprintf(PrefixCacheLease, data.EmailRequestID)
				ok, err := e.redisPort.SetLock(ctx, leaseKey, "1", LeaseDuration)
				if err != nil {
					log.Error(ctx, fmt.Sprintf("Failed to acquire lock for email request %d: %v", data.EmailRequestID, err))
					continue
				}
				if !ok {
					log.Warn(ctx, fmt.Sprintf("Lock already exists for email request %d, skipping", data.EmailRequestID))
					continue
				}

				// 2. Send email
				sendErr := e.emailProviderPort.Send(ctx, data.Provider, data)
				status := constant.EmailSendingStatusSent
				var errMessage string
				timeNow := time.Now()
				sentAt := utils.ToUnixTimeToPointer(&timeNow)

				// 3. Handle 401 (refresh token logic) PER PROVIDER
				if errors.Is(sendErr, common.ErrUnauthorized) {
					log.Warn(ctx, fmt.Sprintf("401 detected for %v. Refreshing token...", data.To))

					providerID := data.Provider.ID // Adjust as per your struct
					// Get or create refreshData for this provider
					v, _ := providerRefreshMap.LoadOrStore(providerID, &refreshData{once: new(sync.Once)})
					refresh := v.(*refreshData)
					refresh.once.Do(func() {
						refreshed, err := e.updateEmailProviderUseCase.UpdateOAuthInfoByRefreshToken(ctx, data.Provider)
						if err != nil {
							refresh.err = err
							log.Error(ctx, "Token refresh failed", err)
							return
						}
						refresh.mu.Lock()
						data.Provider = refreshed
						refresh.mu.Unlock()
					})

					if refresh.err == nil {
						sendErr = e.emailProviderPort.Send(ctx, data.Provider, data)
					}
				}

				if sendErr != nil {
					status = constant.EmailSendingStatusFailed
					errMessage = sendErr.Error()
					log.Error(ctx, fmt.Sprintf("Failed to send email to %v", data.To), sendErr)
				}

				emailRequest := mapEmailRequests[data.EmailRequestID]
				emailRequest.Status = status
				emailRequest.ErrorMessage = errMessage
				emailRequest.SentAt = sentAt

				// 4. Collect results via channel (thread safe)
				results <- emailRequest

				// 5. Release the lease
				if err := e.redisPort.DeleteKey(ctx, leaseKey); err != nil {
					log.Error(ctx, fmt.Sprintf("Failed to release lock for email request %d: %v", data.EmailRequestID, err))
				} else {
					log.Info(ctx, fmt.Sprintf("Released lock for email request %d", data.EmailRequestID))
				}
			}
		}()
	}

	for _, data := range dataSendings {
		jobs <- data
	}
	close(jobs)
	wg.Wait()
	close(results)

	// Collect results
	emailRequestsNoSkip := make([]*entity.EmailRequestEntity, 0, len(dataSendings))
	for r := range results {
		emailRequestsNoSkip = append(emailRequestsNoSkip, r)
	}

	return e.SyncToDB(ctx, emailRequestsNoSkip)
}

func (e EmailSendingUsecase) buildTemplateAndProviderMap(ctx context.Context, emailRequests []*entity.EmailRequestEntity) (map[int64]*entity.EmailProviderEntity, map[int64]*entity.EmailTemplateEntity, error) {
	providerIDs := make([]int64, 0)
	templateIDs := make([]int64, 0)
	for _, emailRequest := range emailRequests {
		templateIDs = append(templateIDs, emailRequest.TemplateId)
		providerIDs = append(providerIDs, emailRequest.EmailProviderID)
	}
	emailProviders, err := e.getEmailProviderUseCase.GetProvidersByIds(ctx, providerIDs)
	if err != nil {
		log.Error(ctx, "Error when get email providers by ids", err)
		return nil, nil, err
	}
	templates, err := e.getEmailTemplateUseCase.GetTemplatesByIDs(ctx, templateIDs)
	if err != nil {
		log.Error(ctx, "Error when get email templates by ids", err)
		return nil, nil, err
	}
	// build email provider map
	emailProviderMap := make(map[int64]*entity.EmailProviderEntity)
	for _, provider := range emailProviders {
		emailProviderMap[provider.ID] = provider
	}
	//build template map
	templateMap := make(map[int64]*entity.EmailTemplateEntity)
	for _, t := range templates {
		templateMap[t.ID] = t
	}
	return emailProviderMap, templateMap, nil
}

func (e EmailSendingUsecase) SyncToDB(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error {
	if _, err := e.updateEmailRequestUsecase.UpdateStatusByBatches(ctx, emailRequests); err != nil {
		log.Error(ctx, "Error when update email requests", err)
		return err
	}
	return nil
}
