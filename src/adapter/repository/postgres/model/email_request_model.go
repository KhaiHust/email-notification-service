package model

import "encoding/json"

type EmailRequestModel struct {
	BaseModel
	TemplateId int64           `gorm:"column:template_id"`
	Recipient  string          `gorm:"column:recipient"`
	Data       json.RawMessage `gorm:"column:data"`
	Status     string          `gorm:"column:status"`
}

func (EmailRequestModel) TableName() string {
	return "email_requests"
}
