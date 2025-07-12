package specification

import (
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"time"
)

type BaseSpecification struct {
	SortOrder     string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
	UpdatedAtFrom *time.Time
	UpdatedAtTo   *time.Time
	Limit         *int64
	Since         *int64
	Until         *int64
}

func ToBaseSpecification(filter *request.BaseFilter) *BaseSpecification {
	if filter == nil {
		return nil
	}
	return &BaseSpecification{
		SortOrder:     filter.SortOrder,
		CreatedAtFrom: utils.FromUnixPointerToTime(filter.CreatedAtFrom),
		CreatedAtTo:   utils.FromUnixPointerToTime(filter.CreatedAtTo),
		UpdatedAtFrom: utils.FromUnixPointerToTime(filter.UpdatedAtFrom),
		UpdatedAtTo:   utils.FromUnixPointerToTime(filter.UpdatedAtTo),
		Limit:         filter.Limit,
		Since:         filter.Since,
		Until:         filter.Until,
	}
}
