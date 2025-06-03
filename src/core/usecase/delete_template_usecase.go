package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IDeleteTemplateUseCase interface {
	DeleteTemplate(ctx context.Context, workspaceID, userID, templateID int64) error
}
type DeleteTemplateUseCase struct {
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort
	workspaceUserRepositoryPort port.IWorkspaceUserRepositoryPort
	databaseTransactionUseCase  IDatabaseTransactionUseCase
}

func (d DeleteTemplateUseCase) DeleteTemplate(ctx context.Context, workspaceID, userID, templateID int64) error {
	workspaceUser, err := d.workspaceUserRepositoryPort.GetWorkspaceUserByWorkspaceIDAndUserID(ctx, workspaceID, userID)
	if err != nil {
		log.Error(ctx, "Failed to get workspace user ", err)
		return err
	}
	if workspaceUser.Role != constant.WorkspaceRoleAdmin {
		log.Error(ctx, "User does not have permission to delete template")
		return common.ErrForbidden
	}
	// Check if the template exists and locked for update
	tx := d.databaseTransactionUseCase.StartTx()
	commitTx := false
	defer func() {
		if r := recover(); r != nil {
			err = common.ErrInternalServer
		}
		if !commitTx || err != nil {
			if errRollback := d.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "[DeleteTemplateUseCase] Error rolling back transaction: %v", errRollback)
			} else {
				log.Error(ctx, "[DeleteTemplateUseCase] Transaction rolled back successfully")
			}
		}
	}()
	emailTemplate, err := d.emailTemplateRepositoryPort.GetTemplateForUpdateByIDAndWorkspaceID(ctx, tx, templateID, workspaceID)
	if err != nil {
		return err
	}
	emailTemplate.LastUpdatedBy = userID
	emailTemplate, err = d.emailTemplateRepositoryPort.DeactivateEmailTemplate(ctx, tx, emailTemplate)
	if err != nil {
		log.Error(ctx, "[DeleteTemplateUseCase] Error deactivating email template: %v", err)
		return err
	}

	if err := d.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "[DeleteTemplateUseCase] Error committing transaction: %v", err)
		return err
	}
	commitTx = true
	return nil
}

func NewDeleteTemplateUseCase(
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort,
	workspaceUserRepositoryPort port.IWorkspaceUserRepositoryPort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
) IDeleteTemplateUseCase {
	return &DeleteTemplateUseCase{
		emailTemplateRepositoryPort: emailTemplateRepositoryPort,
		workspaceUserRepositoryPort: workspaceUserRepositoryPort,
		databaseTransactionUseCase:  databaseTransactionUseCase,
	}
}
