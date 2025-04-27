package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IUpdateEmailRequestUsecase interface {
	UpdateStatusByBatches(ctx context.Context, emailRequestEntity []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateEmailRequestByID(ctx context.Context, emailRequestEntity *entity.EmailRequestEntity) (*entity.EmailRequestEntity, error)
}
type UpdateEmailRequestUsecase struct {
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
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
	//TODO implement me
	panic("implement me")
}

func NewUpdateEmailRequestUsecase(
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
) IUpdateEmailRequestUsecase {
	return &UpdateEmailRequestUsecase{
		emailRequestRepositoryPort: emailRequestRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
	}
}
