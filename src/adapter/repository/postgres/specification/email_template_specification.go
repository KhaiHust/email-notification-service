package specification

import (
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type EmailTemplateSpecification struct {
	Name          *string
	Limit         *int64
	Since         *int64
	Until         *int64
	WorkspaceID   *int64
	DirectTo      string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
	UpdatedAtFrom *time.Time
	UpdatedAtTo   *time.Time
}

func NewEmailTemplateSpecificationQuery(sp *EmailTemplateSpecification) (string, []interface{}, error) {
	query := sq.Select("id, name, workspace_id, created_at, updated_at").From("email_templates").Where("1=1")
	if sp != nil {
		query = addCondition(sp, query)
		// Handle ID filtering
		if sp.DirectTo == constant.ASC {
			if sp.Since != nil {
				query = query.Where(sq.GtOrEq{"id": *sp.Since})
			}
			if sp.Until != nil {
				query = query.Where(sq.LtOrEq{"id": *sp.Until})
			}
			query = query.OrderBy("id ASC")
		} else if sp.DirectTo == constant.DESC {
			if sp.Since != nil {
				query = query.Where(sq.LtOrEq{"id": *sp.Since})
			}
			if sp.Until != nil {
				query = query.Where(sq.GtOrEq{"id": *sp.Until})
			}
			query = query.OrderBy("id DESC")
		}

		if sp.Limit != nil {
			query = query.Limit(uint64(*sp.Limit))
		}
	}
	return query.ToSql()
}
func NewEmailTemplateSpecificationQueryWithCount(sp *EmailTemplateSpecification) (string, []interface{}, error) {
	query := sq.Select("COUNT(*)").From("email_templates").Where("1=1")
	if sp != nil {
		query = addCondition(sp, query)
	}
	return query.ToSql()
}

func addCondition(sp *EmailTemplateSpecification, query sq.SelectBuilder) sq.SelectBuilder {
	if sp.WorkspaceID != nil {
		query = query.Where(sq.Eq{"workspace_id": *sp.WorkspaceID})
	}
	if sp.Name != nil {
		query = query.Where(sq.Like{"name": "%" + *sp.Name + "%"})
	}
	if sp.CreatedAtFrom != nil {
		query = query.Where(sq.GtOrEq{"created_at": *sp.CreatedAtFrom})
	}
	if sp.CreatedAtTo != nil {
		query = query.Where(sq.LtOrEq{"created_at": *sp.CreatedAtTo})
	}
	if sp.UpdatedAtFrom != nil {
		query = query.Where(sq.GtOrEq{"updated_at": *sp.UpdatedAtFrom})
	}
	if sp.UpdatedAtTo != nil {
		query = query.Where(sq.LtOrEq{"updated_at": *sp.UpdatedAtTo})
	}
	return query
}
func ToEmailTemplateSpecification(filter *request.GetListEmailTemplateFilter) *EmailTemplateSpecification {
	if filter == nil {
		return nil
	}
	return &EmailTemplateSpecification{
		Name:          filter.Name,
		Limit:         filter.Limit,
		Since:         filter.Since,
		Until:         filter.Until,
		DirectTo:      filter.DirectTo,
		WorkspaceID:   filter.WorkspaceID,
		CreatedAtFrom: utils.FromUnixPointerToTime(filter.CreatedAtFrom),
		CreatedAtTo:   utils.FromUnixPointerToTime(filter.CreatedAtTo),
		UpdatedAtFrom: utils.FromUnixPointerToTime(filter.UpdatedAtFrom),
		UpdatedAtTo:   utils.FromUnixPointerToTime(filter.UpdatedAtTo),
	}
}
