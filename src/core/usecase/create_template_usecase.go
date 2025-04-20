package usecase

import (
	"context"
	"encoding/json"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
)

type ICreateTemplateUseCase interface {
	CreateTemplate(ctx context.Context, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error)
}
type CreateTemplateUseCase struct {
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort
	databaseTransactionUseCase  IDatabaseTransactionUseCase
}

func (c CreateTemplateUseCase) CreateTemplate(ctx context.Context, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error) {
	extractVariables := utils.ExtractVariablesBySection(template.Subject, template.Body)
	jsonBytes, err := json.Marshal(extractVariables)
	if err != nil {
		log.Error(ctx, "Error when marshal variables to json", err)
		return nil, err
	}
	template.Variables = jsonBytes
	tx := c.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if errRollback := c.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
			log.Error(ctx, "Error when rollback transaction", errRollback)
		} else {
			log.Info(ctx, "Rollback transaction successfully")
		}
	}()
	template, err = c.emailTemplateRepositoryPort.SaveNewTemplate(ctx, tx, template)
	if err != nil {
		log.Error(ctx, "Error when save new template", err)
		return nil, err
	}
	if errCommit := c.databaseTransactionUseCase.CommitTx(tx); errCommit != nil {
		log.Error(ctx, "Error when commit transaction", errCommit)
		return nil, errCommit
	}
	return template, nil
}

func NewCreateTemplateUseCase(emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort, IDatabaseTransactionUseCase IDatabaseTransactionUseCase) ICreateTemplateUseCase {
	return &CreateTemplateUseCase{
		emailTemplateRepositoryPort: emailTemplateRepositoryPort,
		databaseTransactionUseCase:  IDatabaseTransactionUseCase,
	}
}
