package model

import "time"

type CountEmailRequestStatusModel struct {
	Status string     `gorm:"column:status"`
	Total  int64      `gorm:"column:total"`
	Period *time.Time `gorm:"colum:period"`
}
type ChartStatModel struct {
	Period time.Time `json:"period" gorm:"column:period"`
	Sent   int64     `json:"sent" gorm:"column:sent"`
	Error  int64     `json:"error" gorm:"column:error"`
	Open   int64     `json:"open" gorm:"column:open"`
}
type TemplateStatModel struct {
	Sent          int64                `json:"sent" gorm:"column:sent"`
	Error         int64                `json:"error" gorm:"column:error"`
	Open          int64                `json:"open" gorm:"column:open"`
	ProviderStats []*ProviderStatModel `json:"provider_stats" gorm:"-"`
}
type ProviderStatModel struct {
	ProviderID   int64  `json:"provider_id" gorm:"column:email_provider_id"`
	ProviderName string `json:"provider_name" gorm:"column:provider_name"`
	Sent         int64  `json:"sent" gorm:"column:sent"`
	Error        int64  `json:"error" gorm:"column:error"`
	Open         int64  `json:"open" gorm:"column:open"`
}
