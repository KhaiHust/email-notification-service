package specification

import (
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	sq "github.com/Masterminds/squirrel"
)

type EmailRequestSpecification struct {
	WorkspaceIDs     []int64
	EmailTemplateIDs []int64
	Statuses         []string
	RequestID        *string
	Recipient        *string
	*BaseSpecification
}

func NewEmailRequestSpecificationForCount(sp *EmailRequestSpecification) (string, []interface{}, error) {
	builder := sq.
		Select("COUNT(*) as total").
		From("email_requests")

	if sp != nil {
		builder = buildEmailRequestSpecConditions(sp, builder)
	}

	// Generate SQL and args
	query, args, err := builder.ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func NewEmailRequestSpecificationForQuery(sp *EmailRequestSpecification) (string, []interface{}, error) {
	builder := sq.
		Select("id", "template_id", "status", "created_at", "updated_at", "email_provider_id", "request_id", "recipient").
		From("email_requests")

	if sp != nil {
		builder = buildEmailRequestSpecConditions(sp, builder)
		if sp.SortOrder == constant.ASC {
			if sp.Since != nil {
				builder = builder.Where(sq.Gt{"id": *sp.Since})
			}
			if sp.Until != nil {
				builder = builder.Where(sq.Lt{"id": *sp.Until})
			}
			builder = builder.OrderBy("id ASC")
		}
		if sp.SortOrder == constant.DESC {
			if sp.Since != nil {
				builder = builder.Where(sq.Lt{"id": *sp.Since})
			}
			if sp.Until != nil {
				builder = builder.Where(sq.Gt{"id": *sp.Until})
			}
			builder = builder.OrderBy("id DESC")
		}
		if sp.Limit != nil {
			builder = builder.Limit(uint64(*sp.Limit))
		}
	}

	// Generate SQL and args
	query, args, err := builder.ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func buildEmailRequestSpecConditions(sp *EmailRequestSpecification, builder sq.SelectBuilder) sq.SelectBuilder {
	if len(sp.WorkspaceIDs) > 0 {
		builder = builder.Where(sq.Eq{"workspace_id": sp.WorkspaceIDs})
	}
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
	if sp.UpdatedAtFrom != nil {
		builder = builder.Where(sq.GtOrEq{"updated_at": *sp.UpdatedAtFrom})
	}
	if sp.UpdatedAtTo != nil {
		builder = builder.Where(sq.LtOrEq{"updated_at": *sp.UpdatedAtTo})
	}
	if sp.RequestID != nil && *sp.RequestID != "" {
		builder = builder.Where(sq.Eq{"request_id": *sp.RequestID})
	}
	if sp.Recipient != nil && *sp.Recipient != "" {
		builder = builder.Where(sq.Eq{"recipient": *sp.Recipient})
	}
	return builder
}
func NewEmailRequestSpecificationForCountStatus(sp *EmailRequestSpecification) (string, []interface{}, error) {
	builder := sq.
		Select("template_id", "status", "COUNT(*) as total").
		From("email_requests").
		GroupBy("template_id", "status")

	if sp != nil {
		builder = buildEmailRequestSpecConditions(sp, builder)
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
		WorkspaceIDs:      filter.WorkspaceIDs,
		EmailTemplateIDs:  filter.EmailTemplateIDs,
		Statuses:          filter.Statuses,
		BaseSpecification: ToBaseSpecification(filter.BaseFilter),
		RequestID:         filter.RequestID,
		Recipient:         filter.Email,
	}
}
