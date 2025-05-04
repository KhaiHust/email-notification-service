package specification

import (
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type EmailRequestSpecification struct {
	EmailTemplateIDs []int64
	Statuses         []string
	CreatedAtFrom    *time.Time
	CreatedAtTo      *time.Time
}

func NewEmailRequestSpecificationForCountStatus(sp *EmailRequestSpecification) (string, []interface{}, error) {
	builder := sq.
		Select("template_id", "status", "COUNT(*) as total").
		From("email_requests").
		GroupBy("template_id", "status")

	if sp != nil {
		if len(sp.EmailTemplateIDs) > 0 {
			builder = builder.Where(sq.Eq{"template_id": sp.EmailTemplateIDs})
		}
		if len(sp.Statuses) > 0 {
			builder = builder.Where(sq.Eq{"status": sp.Statuses})
		}
		if sp.CreatedAtFrom != nil {
			builder = builder.Where(sq.GtOrEq{"created_at": *sp.CreatedAtFrom})
		}
		if sp.CreatedAtTo != nil {
			builder = builder.Where(sq.LtOrEq{"created_at": *sp.CreatedAtTo})
		}
	}

	// Generate SQL and args
	query, args, err := builder.ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}
func ToEmailRequestSpecification(filter *request.EmailRequestFilter) *EmailRequestSpecification {
	if filter == nil {
		return nil
	}
	return &EmailRequestSpecification{
		EmailTemplateIDs: filter.EmailTemplateIDs,
		Statuses:         filter.Statuses,
		CreatedAtFrom:    utils.FromUnixPointerToTime(filter.CreatedAtFrom),
		CreatedAtTo:      utils.FromUnixPointerToTime(filter.CreatedAtTo),
	}
}
