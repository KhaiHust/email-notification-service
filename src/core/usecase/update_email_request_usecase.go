package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"time"
)

type IUpdateEmailRequestUsecase interface {
	UpdateStatusByBatches(ctx context.Context, emailRequestEntity []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateEmailRequestByID(ctx context.Context, emailRequestEntity *entity.EmailRequestEntity) (*entity.EmailRequestEntity, error)
}
type UpdateEmailRequestUsecase struct {
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
	emailLogRepositoryPort     port.IEmailLogRepositoryPort
}

func (u UpdateEmailRequestUsecase) UpdateEmailRequestByID(ctx context.Context, emailRequestEntity *entity.EmailRequestEntity) (*entity.EmailRequestEntity, error) {
	var err error
	tx := u.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := u.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Error when rollback transaction", errRollback)
			} else {
				err = u.databaseTransactionUseCase.CommitTx(tx)
			}
		}
	}()
	emailRequestEntity, err = u.emailRequestRepositoryPort.UpdateEmailRequestByID(ctx, tx, emailRequestEntity)
	if err != nil {
		log.Error(ctx, "Error when update email request by id", err)
		return nil, err
	}
	if errCommit := u.databaseTransactionUseCase.CommitTx(tx); errCommit != nil {
		log.Error(ctx, "Error when commit transaction", errCommit)
		return nil, errCommit
	}
	return emailRequestEntity, nil
}

func (u UpdateEmailRequestUsecase) UpdateStatusByBatches(ctx context.Context, emailRequestEntity []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	var err error
	tx := u.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := u.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Error when rollback transaction", errRollback)
			} else {
				log.Info(ctx, "Transaction rolled back successfully")
			}
		}
	}()
	for _, emailRequest := range emailRequestEntity {
		emailRequest, err = u.emailRequestRepositoryPort.UpdateEmailRequestByID(ctx, tx, emailRequest)
		if err != nil {
			log.Error(ctx, "Error when update email request by id", err)
			return nil, err
		}
		emailLog := u.toEmailLogEntity(emailRequest)
		if _, err = u.emailLogRepositoryPort.SaveNewEmailLog(ctx, tx, emailLog); err != nil {
			log.Error(ctx, "Error when save email log", err)
			return nil, err
		}
	}
	if errCommit := u.databaseTransactionUseCase.CommitTx(tx); errCommit != nil {
		log.Error(ctx, "Error when commit transaction", errCommit)
		return nil, errCommit
	}
	return emailRequestEntity, nil
}

func NewUpdateEmailRequestUsecase(
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	emailLogRepositoryPort port.IEmailLogRepositoryPort,
) IUpdateEmailRequestUsecase {
	return &UpdateEmailRequestUsecase{
		emailRequestRepositoryPort: emailRequestRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
		emailLogRepositoryPort:     emailLogRepositoryPort,
	}
}
func (u UpdateEmailRequestUsecase) toEmailLogEntity(emailRequest *entity.EmailRequestEntity) *entity.EmailLogsEntity {
	var loggedAt int64
	if emailRequest.Status == constant.EmailSendingStatusSent {
		loggedAt = *emailRequest.SentAt
	}
	if emailRequest.Status == constant.EmailSendingStatusOpened {
		loggedAt = *emailRequest.OpenedAt
	}
	if emailRequest.Status == constant.EmailSendingStatusFailed {
		sendAt := emailRequest.SentAt
		if sendAt == nil {
			sendAt = utils.ToInt64Pointer(time.Now().Unix())
		}
		loggedAt = *sendAt
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
