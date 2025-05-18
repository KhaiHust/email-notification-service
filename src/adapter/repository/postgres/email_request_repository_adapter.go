package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/specification"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmailRequestRepositoryAdapter struct {
	base
}

func (e EmailRequestRepositoryAdapter) GetEmailRequestForUpdateByIDOrTrackingID(ctx context.Context, tx *gorm.DB, emailRequestID int64, trackingID string) (*entity.EmailRequestEntity, error) {
	emailRequestModel := &model.EmailRequestModel{}
	if err := tx.WithContext(ctx).Model(&model.EmailRequestModel{}).Clauses(
		clause.Locking{Strength: clause.LockingStrengthUpdate},
	).
		Where("id = ? OR tracking_id = ?", emailRequestID, trackingID).First(emailRequestModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailRequestEntity(emailRequestModel), nil
}

func (e EmailRequestRepositoryAdapter) CountAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) (int64, error) {
	emailRequestSpec := specification.ToEmailRequestSpecification(filter)
	query, args, err := specification.NewEmailRequestSpecificationForCount(emailRequestSpec)
	if err != nil {
		return 0, err
	}
	var count int64
	if err := e.db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (e EmailRequestRepositoryAdapter) GetAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestEntity, error) {
	emailRequestSpec := specification.ToEmailRequestSpecification(filter)
	query, args, err := specification.NewEmailRequestSpecificationForQuery(emailRequestSpec)
	if err != nil {
		return nil, err
	}
	var emailRequestModels []*model.EmailRequestModel
	if err := e.db.WithContext(ctx).Preload("EmailTemplateModel", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).
		Preload("EmailProviderModel", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, provider")
		}).
		Raw(query, args...).Find(&emailRequestModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailRequestEntity(emailRequestModels), nil
}

func (e EmailRequestRepositoryAdapter) CountEmailRequestStatuses(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestStatusCountEntity, error) {
	emailRequestSpec := specification.ToEmailRequestSpecification(filter)
	query, args, err := specification.NewEmailRequestSpecificationForCountStatus(emailRequestSpec)
	if err != nil {
		return nil, err
	}
	var emailRequestStatusCountModels []*model.EmailRequestStatusCountModel
	if err := e.db.WithContext(ctx).Raw(query, args...).Scan(&emailRequestStatusCountModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailStatusCountEntity(emailRequestStatusCountModels), nil
}

func (e EmailRequestRepositoryAdapter) GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error) {
	emailRequestModel := &model.EmailRequestModel{}
	if err := e.db.WithContext(ctx).Model(&model.EmailRequestModel{}).Where("id = ?", emailRequestID).First(emailRequestModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailRequestEntity(emailRequestModel), nil
}

func (e EmailRequestRepositoryAdapter) UpdateEmailRequestByID(ctx context.Context, tx *gorm.DB, emailRequest *entity.EmailRequestEntity) (*entity.EmailRequestEntity, error) {
	emailRequestModel := mapper.ToEmailRequestModel(emailRequest)
	if err := tx.WithContext(ctx).Model(&model.EmailRequestModel{}).Where("id = ?", emailRequestModel.ID).Updates(emailRequestModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailRequestEntity(emailRequestModel), nil
}

func (e EmailRequestRepositoryAdapter) UpdateStatusByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (e EmailRequestRepositoryAdapter) SaveEmailRequestByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error) {
	emailRequestModels := mapper.ToListEmailRequestModel(emailRequests)
	if err := tx.WithContext(ctx).Model(&model.EmailRequestModel{}).Create(emailRequestModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailRequestEntity(emailRequestModels), nil
}

func NewEmailRequestRepositoryAdapter(db *gorm.DB) port.IEmailRequestRepositoryPort {
	return &EmailRequestRepositoryAdapter{
		base: base{
			db: db,
		}}
}
