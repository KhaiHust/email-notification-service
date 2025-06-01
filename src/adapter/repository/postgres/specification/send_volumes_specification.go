package specification

import (
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type SendVolumeSpecification struct {
	WorkspaceID int64
	StartDate   *time.Time
	EndDate     *time.Time
}

func NewSendVolumeSpecification(filter *request.SendVolumeFilter) *SendVolumeSpecification {
	return &SendVolumeSpecification{
		WorkspaceID: filter.WorkspaceId,
		StartDate:   utils.FromUnixPointerToTime(filter.StartDate),
		EndDate:     utils.FromUnixPointerToTime(filter.EndDate),
	}
}

func (s *SendVolumeSpecification) ToSendVolumeQuery() (string, []interface{}, error) {
	builder := sq.
		Select("DATE(created_at) AS date", "COUNT(*) AS total").
		From("email_requests").
		Where(sq.Eq{"workspace_id": s.WorkspaceID}).
		Where(sq.And{
			sq.GtOrEq{"created_at": s.StartDate},
			sq.LtOrEq{"created_at": s.EndDate},
		}).
		GroupBy("DATE(created_at)").
		OrderBy("DATE(created_at) ASC")

	query, args, err := builder.ToSql()
	return query, args, err
}
func (s *SendVolumeSpecification) ToSendVolumeQueryByProvider() (string, []interface{}, error) {
	builder := sq.
		Select("email_provider_id AS provider_id", "DATE(created_at) AS date", "COUNT(*) AS total").
		From("email_requests").
		Where(sq.Eq{"workspace_id": s.WorkspaceID}).
		Where(sq.And{
			sq.GtOrEq{"created_at": s.StartDate},
			sq.LtOrEq{"created_at": s.EndDate},
		}).
		GroupBy("email_provider_id", "DATE(created_at)").
		OrderBy("email_provider_id ASC", "DATE(created_at) ASC")

	query, args, err := builder.ToSql()
	return query, args, err
}
func (s *SendVolumeSpecification) ToSendVolumeByProviderQuery() (string, []interface{}, error) {
	// count total and total error group by provider
	builder := sq.
		Select("email_provider_id AS provider_id", "COUNT(*) AS total",
			"SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) AS total_error",
			"SUM(CASE WHEN status IN ('SENT','OPENED') THEN 1 ELSE 0 END) AS total_sent").
		From("email_requests").
		Where(sq.Eq{"workspace_id": s.WorkspaceID}).
		Where(sq.And{
			sq.GtOrEq{"created_at": s.StartDate},
			sq.LtOrEq{"created_at": s.EndDate},
		}).
		GroupBy("email_provider_id").
		OrderBy("email_provider_id ASC")
	return builder.ToSql()
}
