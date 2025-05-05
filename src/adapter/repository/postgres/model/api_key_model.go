package model

import "time"

type ApiKeyModel struct {
	BaseModel
	WorkspaceID int64      `gorm:"column:workspace_id"`
	Name        string     `gorm:"column:name"`
	KeyHash     string     `gorm:"column:key_hash"`
	RawPrefix   string     `gorm:"column:raw_prefix"`
	Environment string     `gorm:"column:environment"`
	ExpiresAt   *time.Time `gorm:"column:expires_at"`
	Revoked     bool       `gorm:"column:revoked"`
}

func (a *ApiKeyModel) TableName() string {
	return "api_keys"
}
