package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"time"
)

type IEmailTrackingUsecase interface {
	OpenEmailTracking(ctx context.Context, encryptTrackingID string) error
}
type EmailTrackingUsecase struct {
	encryptUseCase             IEncryptUseCase
	updateEmailRequestUsecase  IUpdateEmailRequestUsecase
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
	emailLogRepositoryPort     port.IEmailLogRepositoryPort
}

func (e EmailTrackingUsecase) OpenEmailTracking(ctx context.Context, encryptTrackingID string) error {
	trackingID, err := e.encryptUseCase.DecryptTrackingID(ctx, encryptTrackingID)
	if err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when decrypt tracking id", err)
		return common.ErrInvalidEmailTrackingID
	}
	go e.SyncToDb(ctx, trackingID)
	return nil
}

func (e EmailTrackingUsecase) SyncToDb(ctx context.Context, trackingID string) {
	var err error
	tx := e.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := e.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "[EmailTrackingUsecase] Error when rollback transaction", errRollback)
			} else {
				log.Info(ctx, "[EmailTrackingUsecase] Rollback transaction successfully")
			}
		}
	}()
	emailRequestEntity, err := e.emailRequestRepositoryPort.GetEmailRequestForUpdateByIDOrTrackingID(ctx, tx, 0, trackingID)
	if err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when get email request by tracking id", err)
		return
	}
	emailRequestEntity.Status = constant.EmailSendingStatusOpened
	emailRequestEntity.OpenedAt = utils.ToInt64Pointer(time.Now().Unix())
	emailRequestEntity, err = e.updateEmailRequestUsecase.UpdateEmailRequestByID(ctx, emailRequestEntity)
	if err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when update email request status to opened", err)
		return
	}
	emailLog := e.toEmailLogEntity(emailRequestEntity)
	if _, err = e.emailLogRepositoryPort.SaveNewEmailLog(ctx, tx, emailLog); err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when save email log", err)
		return
	}
	if err = e.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when commit transaction", err)
		return
	}
}
func (e EmailTrackingUsecase) toEmailLogEntity(emailRequest *entity.EmailRequestEntity) *entity.EmailLogsEntity {
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
func NewEmailTrackingUsecase(
	encryptUseCase IEncryptUseCase,
	updateEmailRequestUsecase IUpdateEmailRequestUsecase,
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	emailLogRepositoryPort port.IEmailLogRepositoryPort,
) IEmailTrackingUsecase {
	return &EmailTrackingUsecase{
		encryptUseCase:             encryptUseCase,
		updateEmailRequestUsecase:  updateEmailRequestUsecase,
		emailRequestRepositoryPort: emailRequestRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
		emailLogRepositoryPort:     emailLogRepositoryPort,
	}
}
