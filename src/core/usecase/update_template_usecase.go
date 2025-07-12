package usecase

import (
	"context"
	"encoding/json"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"strconv"
	"time"
)

type IUpdateEmailTemplateUseCase interface {
	UpdateEmailTemplate(ctx context.Context, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error)
}
type UpdateEmailTemplateUseCase struct {
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort
	encryptUseCase              IEncryptUseCase
	databaseTransactionUseCase  IDatabaseTransactionUseCase
}

func (u UpdateEmailTemplateUseCase) UpdateEmailTemplate(ctx context.Context, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error) {
	//validate version
	version, err := u.encryptUseCase.DecryptVersionTemplate(ctx, template.Version)
	if err != nil {
		log.Error(ctx, "[UpdateEmailTemplateUseCase] Error decrypting version: %v", err)
		return nil, common.ErrEmailTemplateVersionInvalidFormat
	}
	extractVariables := utils.ExtractVariablesBySection(template.Subject, template.Body)
	jsonBytes, err := json.Marshal(extractVariables)
	if err != nil {
		log.Error(ctx, "Error when marshal variables to json", err)
		return nil, err
	}
	template.Variables = jsonBytes
	tx := u.databaseTransactionUseCase.StartTx()
	commitTx := false
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if !commitTx || err != nil {
			if errRollback := u.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "[UpdateEmailTemplateUseCase] Error rolling back transaction: %v", errRollback)
			} else {
				log.Error(ctx, "[UpdateEmailTemplateUseCase] Transaction rolled back successfully")
			}
		}
	}()

	//get template for update
	emailTemplate, err := u.emailTemplateRepositoryPort.GetTemplateForUpdateByIDAndWorkspaceID(ctx, tx, template.ID, template.WorkspaceId)
	if err != nil {
		log.Error(ctx, "[UpdateEmailTemplateUseCase] Error getting template for update: %v", err)
		return nil, err
	}
	//check version
	if emailTemplate.Version != version {
		log.Error(ctx, "[UpdateEmailTemplateUseCase] Version template not match")
		return emailTemplate, common.ErrEmailTemplateVersionNotMatch
	}
	newVersion := time.Now().Unix()
	template.Version = strconv.FormatInt(newVersion, 10)
	//save template
	emailTemplate, err = u.emailTemplateRepositoryPort.UpdateTemplate(ctx, tx, template)
	if err != nil {
		log.Error(ctx, "[UpdateEmailTemplateUseCase] Error saving template: %v", err)
		return nil, err
	}
	//commit transaction
	if err := u.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "[UpdateEmailTemplateUseCase] Error committing transaction: %v", err)
		return nil, err
	}
	commitTx = true
	return emailTemplate, nil

}

func NewUpdateEmailTemplateUseCase(
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	encryptUseCase IEncryptUseCase,
) IUpdateEmailTemplateUseCase {
	return &UpdateEmailTemplateUseCase{
		emailTemplateRepositoryPort: emailTemplateRepositoryPort,
		databaseTransactionUseCase:  databaseTransactionUseCase,
		encryptUseCase:              encryptUseCase,
	}
}
