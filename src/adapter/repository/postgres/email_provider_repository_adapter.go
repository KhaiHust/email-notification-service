package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailProviderRepositoryAdapter struct {
	base
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) (*entity.EmailProviderEntity, error) {
	//build raw sql
	sql := "SELECT ep.id, ep.workspace_id, ep.provider, ep.email, ep.from_name, ep.environment FROM email_providers ep " +
		"JOIN workspaces w ON ep.workspace_id = w.id " +
		"WHERE w.code = ? AND ep.provider = ?"
	var emailProviderModel model.EmailProviderModel
	if err := e.db.WithContext(ctx).Raw(sql, workspaceCode, provider).Scan(&emailProviderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailProviderEntity(&emailProviderModel), nil
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByWorkspaceID(ctx context.Context, workspaceID int64) (*entity.EmailProviderEntity, error) {
	var emailProviderModels *model.EmailProviderModel
	if err := e.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).First(&emailProviderModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModels), nil
}

func (e EmailProviderRepositoryAdapter) UpdateEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	emailProviderModel := mapper.ToEmailProviderModel(emailProvider)
	if err := tx.WithContext(ctx).Model(&model.EmailProviderModel{}).
		Where("id = ?", emailProvider.ID).
		Updates(emailProviderModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModel), nil
}

func (e EmailProviderRepositoryAdapter) GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error) {
	var emailProviderModel model.EmailProviderModel
	if err := e.db.WithContext(ctx).Where("id = ?", ID).First(&emailProviderModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailProviderEntity(&emailProviderModel), nil
}

func (e EmailProviderRepositoryAdapter) SaveEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	emailProviderModel := mapper.ToEmailProviderModel(emailProvider)

	if err := tx.WithContext(ctx).Create(emailProviderModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModel), nil
}

func NewEmailProviderRepositoryAdapter(db *gorm.DB) port.IEmailProviderRepositoryPort {
	return &EmailProviderRepositoryAdapter{
		base: base{db},
	}
}
