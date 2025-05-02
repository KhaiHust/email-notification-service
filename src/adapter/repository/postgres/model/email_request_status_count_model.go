package model

type EmailRequestStatusCountModel struct {
	EmailTemplateID int64  `gorm:"column:template_id"`
	Status          string `gorm:"column:status"`
	Total           int64  `gorm:"column:total"`
}
