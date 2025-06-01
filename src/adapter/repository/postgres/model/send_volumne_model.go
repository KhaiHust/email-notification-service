package model

import "time"

type SendVolumeProviderByDateModel struct {
	ProviderID int64     `gorm:"column:provider_id"`
	Total      int64     `gorm:"column:total"`
	Date       time.Time `gorm:"column:date"`
}
type SendVolumeByModel struct {
	Total int64     `gorm:"column:total"`
	Date  time.Time `gorm:"column:date"`
}
type SendVolumeByProviderModel struct {
	ProviderID int64 `gorm:"column:provider_id"`
	Total      int64 `gorm:"column:total"`
	TotalError int64 `gorm:"column:total_error"`
	TotalSent  int64 `gorm:"column:total_sent"`
}
