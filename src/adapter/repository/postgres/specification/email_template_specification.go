package specification

import (
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type EmailTemplateSpecification struct {
	Name        *string
	WorkspaceID *int64
	*BaseSpecification
}

func NewEmailTemplateSpecificationQuery(sp *EmailTemplateSpecification) (string, []interface{}, error) {
	query := sq.Select("id, name, workspace_id, created_at, updated_at").From("email_templates").Where("1=1 AND active IS TRUE")
	if sp != nil {
		query = addCondition(sp, query)
		// Handle ID filtering
		if sp.SortOrder == constant.ASC {
			if sp.Since != nil {
				query = query.Where(sq.GtOrEq{"id": *sp.Since})
			}
			if sp.Until != nil {
				query = query.Where(sq.LtOrEq{"id": *sp.Until})
			}
			query = query.OrderBy("id ASC")
		} else if sp.SortOrder == constant.DESC {
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
	query := sq.Select("COUNT(*)").From("email_templates").Where("1=1 AND active IS TRUE")
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
		Name:              filter.Name,
		BaseSpecification: ToBaseSpecification(filter.BaseFilter),
		WorkspaceID:       filter.WorkspaceID,
	}
}

var validIntervals = map[string]bool{"day": true, "week": true, "month": true}

func BuildChartStatsQuery(filter *request.TemplateMetricFilter) (string, []interface{}, error) {
	if !validIntervals[filter.Interval] {
		return "", nil, fmt.Errorf("invalid interval: %s", filter.Interval)
	}
	timeZone, _ := time.Now().Zone()

	periodExpr := fmt.Sprintf(
		"DATE_TRUNC('%s', er.created_at AT TIME ZONE 'UTC' AT TIME ZONE '%s')",
		filter.Interval, timeZone,
	)

	// Subquery for actual data
	subBuilder := sq.
		Select(
			periodExpr+" AS period",
			"COUNT(*) FILTER (WHERE status = 'SENT') AS sent",
			"COUNT(*) FILTER (WHERE status = 'ERROR') AS error",
			"COUNT(*) FILTER (WHERE status = 'OPENED') AS open",
			"COUNT(*) FILTER (WHERE status = 'SCHEDULED') AS scheduled",
		).
		From("email_requests er").
		Where(sq.Eq{"er.workspace_id": filter.WorkspaceID, "er.template_id": filter.TemplateID}).
		Where(sq.GtOrEq{"er.created_at": utils.FromUnixPointerToTime(filter.StartDate)}).
		Where(sq.LtOrEq{"er.created_at": utils.FromUnixPointerToTime(filter.EndDate)}).
		GroupBy("period")

	subQuery, args, err := subBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Format for generate_series
	// For 'day' -> INTERVAL '1 day'
	// For 'week' -> INTERVAL '1 week'
	// For 'month' -> INTERVAL '1 month'
	intervalMap := map[string]string{
		"day":   "1 day",
		"week":  "1 week",
		"month": "1 month",
	}
	intervalStr, ok := intervalMap[filter.Interval]
	if !ok {
		return "", nil, fmt.Errorf("unsupported interval: %s", filter.Interval)
	}

	// Final SQL (raw, with args for the subquery)
	sql := fmt.Sprintf(`
WITH periods AS (
  SELECT generate_series(
    DATE_TRUNC('%[1]s', ? AT TIME ZONE 'UTC' AT TIME ZONE '%[2]s'),
    DATE_TRUNC('%[1]s', ? AT TIME ZONE 'UTC' AT TIME ZONE '%[2]s'),
    INTERVAL '%[3]s'
  ) AS period
)
SELECT
  periods.period,
  COALESCE(stats.sent, 0) AS sent,
  COALESCE(stats.error, 0) AS error,
  COALESCE(stats.open, 0) AS open,
  COALESCE(stats.scheduled, 0) AS scheduled
FROM periods
LEFT JOIN (%s) AS stats
  ON stats.period = periods.period::timestamptz
ORDER BY periods.period ASC
`, filter.Interval, timeZone, intervalStr, subQuery)

	// The $1 and $2 will be StartDate and EndDate (as time.Time)
	args = append([]interface{}{
		utils.FromUnixPointerToTime(filter.StartDate),
		utils.FromUnixPointerToTime(filter.EndDate),
	}, args...)

	return sql, args, nil
}
func BuildTemplateStatQuery(filter *request.TemplateMetricFilter) (string, []interface{}, error) {
	builder := sq.
		Select(
			"COUNT(*) FILTER (WHERE status = 'SENT') AS sent",
			"COUNT(*) FILTER (WHERE status = 'ERROR') AS error",
			"COUNT(*) FILTER (WHERE status = 'OPENED') AS open",
		).
		From("email_requests").
		Where(sq.Eq{"workspace_id": filter.WorkspaceID, "template_id": filter.TemplateID}).
		Where(sq.GtOrEq{"created_at": utils.FromUnixPointerToTime(filter.StartDate)}).
		Where(sq.LtOrEq{"created_at": utils.FromUnixPointerToTime(filter.EndDate)})

	return builder.ToSql()
}
func BuildProviderStatQuery(filter *request.TemplateMetricFilter) (string, []interface{}, error) {
	builder := sq.
		Select(
			"email_provider_id",
			"ep.provider AS provider_name",
			"COUNT(*) FILTER (WHERE er.status = 'SENT') AS sent",
			"COUNT(*) FILTER (WHERE er.status = 'ERROR') AS error",
			"COUNT(*) FILTER (WHERE er.status = 'OPENED') AS open",
		).
		From("email_requests er").
		Join("email_providers ep ON er.email_provider_id = ep.id").
		Where(sq.Eq{"er.workspace_id": filter.WorkspaceID, "er.template_id": filter.TemplateID}).
		Where(sq.GtOrEq{"er.created_at": utils.FromUnixPointerToTime(filter.StartDate)}).
		Where(sq.LtOrEq{"er.created_at": utils.FromUnixPointerToTime(filter.EndDate)}).
		GroupBy("email_provider_id", "ep.provider").
		OrderBy("provider_name ASC")

	return builder.ToSql()
}
