package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type ICreateEmailRequestUsecase interface {
	CreateEmailRequests(ctx context.Context, emailRequestEntities []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
}
type CreateEmailRequestUsecase struct {
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
}

func (c CreateEmailRequestUsecase) CreateEmailRequests(ctx context.Context, emailRequestEntities []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	var err error
	tx := c.databaseTransactionUseCase.StartTx()
	commitTx := false
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if !commitTx || err != nil {
			if errRollback := c.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Error when rollback transaction", errRollback)
			} else {
				log.Info(ctx, "Rollback transaction successfully")
			}
		}
	}()
	emailRequestEntities, err = c.emailRequestRepositoryPort.SaveEmailRequestByBatches(ctx, tx, emailRequestEntities)
	if err != nil {
		log.Error(ctx, "Error when save email request", err)
		return nil, err
	}
	if err = c.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Error when commit transaction", err)
		return nil, err
	}
	commitTx = true
	return emailRequestEntities, nil
}

func NewCreateEmailRequestUsecase(
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
) ICreateEmailRequestUsecase {
	return &CreateEmailRequestUsecase{
		emailRequestRepositoryPort: emailRequestRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
	}
}
