package specification

import (
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"strings"
	"time"
)

type EmailRequestSpecification struct {
	EmailTemplateIDs []int64
	Statuses         []string
	CreatedAtFrom    *time.Time
	CreatedAtTo      *time.Time
}

func NewEmailRequestSpecificationForCountStatus(sp *EmailRequestSpecification) (string, []interface{}, error) {
	query := `
		SELECT template_id, status, COUNT(*) as total
		FROM email_requests
		WHERE 1=1
	`

	var args []interface{}

	if sp != nil {
		if len(sp.EmailTemplateIDs) > 0 {
			query += " AND template_id IN (?" + strings.Repeat(",?", len(sp.EmailTemplateIDs)-1) + ")"
			args = append(args, sp.EmailTemplateIDs)
		}
		if len(sp.Statuses) > 0 {
			query += " AND status IN (?" + strings.Repeat(",?", len(sp.Statuses)-1) + ")"
			args = append(args, sp.Statuses)
		}
		if sp.CreatedAtFrom != nil {
			query += " AND created_at >= ?"
			args = append(args, *sp.CreatedAtFrom)
		}
		if sp.CreatedAtTo != nil {
			query += " AND created_at <= ?"
			args = append(args, *sp.CreatedAtTo)
		}
	}

	query += " GROUP BY template_id, status"

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
