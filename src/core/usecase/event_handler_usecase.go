package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"sync"
	"time"
)

type IEventHandlerUsecase interface {
	SendEmailRequestHandler(ctx context.Context, req *request.EmailSendingRequestDto) error
	SyncEmailRequestHandler(ctx context.Context, emailRequest *entity.EmailRequestEntity) error
}
type EventHandlerUsecase struct {
	emailSendingUsecase        IEmailSendingUsecase
	updateEmailRequestUsecase  IUpdateEmailRequestUsecase
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
	emailLogRepositoryPort     port.IEmailLogRepositoryPort
	scheduleEmailUsecase       IScheduleEmailUsecase
}

func (e EventHandlerUsecase) SyncEmailRequestHandler(ctx context.Context, emailRequest *entity.EmailRequestEntity) error {
	tx := e.databaseTransactionUseCase.StartTx()
	commit := false
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if !commit || err != nil {
			if err := e.databaseTransactionUseCase.RollbackTx(tx); err != nil {
				log.Error(ctx, "Error when rollback transaction", err)
			} else {
				log.Info(ctx, "Rollback transaction success")
			}
		}
	}()
	emailRequestEntity, err := e.emailRequestRepositoryPort.GetEmailRequestForUpdateByIDOrTrackingID(ctx, tx, emailRequest.ID, emailRequest.TrackingID)
	if err != nil {
		log.Error(ctx, "Error when get email request by id", err)
		return err
	}
	emailRequestEntity.Status = emailRequest.Status
	emailRequestEntity.ErrorMessage = emailRequest.ErrorMessage
	emailRequestEntity.SentAt = emailRequest.SentAt
	if emailRequest.Status == constant.EmailSendingStatusOpened {
		emailRequestEntity.OpenedAt = emailRequest.OpenedAt
		emailRequestEntity.OpenedCount += 1
	}
	if _, err = e.emailRequestRepositoryPort.UpdateEmailRequestByID(ctx, tx, emailRequestEntity); err != nil {
		log.Error(ctx, "Error when update email request by id", err)
		return err
	}
	//add to log table
	emailLog := e.toEmailLogEntity(emailRequestEntity)
	if _, err = e.emailLogRepositoryPort.SaveNewEmailLog(ctx, tx, emailLog); err != nil {
		log.Error(ctx, "Error when save email log", err)
		return err
	}
	if err = e.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Error when commit transaction", err)
		return err
	}
	commit = true
	return nil
}

func (e EventHandlerUsecase) SendEmailRequestHandler(ctx context.Context, req *request.EmailSendingRequestDto) error {
	//seperate schedule and normal email request
	emailRequestIDs := make([]int64, 0, len(req.Datas))
	mapData := make(map[int64]*request.EmailSendingData)
	for _, emailRequest := range req.Datas {
		emailRequestIDs = append(emailRequestIDs, emailRequest.EmailRequestID)
		mapData[emailRequest.EmailRequestID] = emailRequest
	}
	emailRequests, err := e.emailRequestRepositoryPort.GetEmailRequestByIDs(ctx, emailRequestIDs)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get email requests by ids: %v", emailRequestIDs), err)
		return err
	}
	if len(emailRequests) == 0 {
		log.Warn(ctx, "No email requests found for the given IDs")
		return nil
	}
	for _, emailRequest := range emailRequests {
		if exist, ok := mapData[emailRequest.ID]; ok {
			if exist.IsRetry {
				emailRequest.IsRetry = true
				emailRequest.RetryCount += 1
			}
		}
	}
	//filter email requests that need to be scheduled
	var wg sync.WaitGroup
	emailRequestsSchedules, emailRequestSending := e.FilterEmailRequestToCreateTasks(emailRequests)

	if len(emailRequestsSchedules) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := e.scheduleEmailUsecase.ScheduleEmails(ctx, emailRequestsSchedules); err != nil {
				log.Error(ctx, "Error when schedule email", err)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.emailSendingUsecase.SendSyncs(ctx, emailRequestSending); err != nil {
			log.Error(ctx, "Error when sending email batches", err)
		}
	}()

	wg.Wait()
	return nil
}
func (e EventHandlerUsecase) toEmailLogEntity(emailRequest *entity.EmailRequestEntity) *entity.EmailLogsEntity {
	var loggedAt int64
	if emailRequest.Status == constant.EmailSendingStatusSent {
		loggedAt = *emailRequest.SentAt
	}
	if emailRequest.Status == constant.EmailSendingStatusOpened {
		loggedAt = *emailRequest.OpenedAt
	}
	if emailRequest.Status == constant.EmailSendingStatusFailed {
		loggedAt = *emailRequest.SentAt
	}
	return &entity.EmailLogsEntity{
		EmailRequestID:  emailRequest.ID,
		Status:          emailRequest.Status,
		ErrorMessage:    emailRequest.ErrorMessage,
		LoggedAt:        loggedAt,
		RetryCount:      emailRequest.RetryCount,
		RequestID:       emailRequest.RequestID,
		WorkspaceID:     emailRequest.WorkspaceID,
		EmailProviderID: emailRequest.EmailProviderID,
		TemplateId:      emailRequest.TemplateId,
		Recipient:       emailRequest.Recipient,
	}

}

func NewEventHandlerUsecase(
	emailSendingUsecase IEmailSendingUsecase,
	updateEmailRequestUsecase IUpdateEmailRequestUsecase,
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	emailLogRepositoryPort port.IEmailLogRepositoryPort,
	scheduleEmailUsecase IScheduleEmailUsecase,
) IEventHandlerUsecase {
	return &EventHandlerUsecase{
		emailSendingUsecase:        emailSendingUsecase,
		updateEmailRequestUsecase:  updateEmailRequestUsecase,
		emailRequestRepositoryPort: emailRequestRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
		emailLogRepositoryPort:     emailLogRepositoryPort,
		scheduleEmailUsecase:       scheduleEmailUsecase,
	}
}
func (e EventHandlerUsecase) FilterEmailRequestToCreateTasks(emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, []*entity.EmailRequestEntity) {
	scheduledEmailRequests := make([]*entity.EmailRequestEntity, 0)
	sendingEmailRequests := make([]*entity.EmailRequestEntity, 0)
	for _, emailRequest := range emailRequests {
		if emailRequest.SendAt != nil && *emailRequest.SendAt > time.Now().Unix() {
			scheduledEmailRequests = append(scheduledEmailRequests, emailRequest)
		} else {
			sendingEmailRequests = append(sendingEmailRequests, emailRequest)
		}
	}
	return scheduledEmailRequests, sendingEmailRequests
}
