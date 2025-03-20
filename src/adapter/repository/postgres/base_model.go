package postgres

import "time"

type BaseModel struct {
	ID       int64     `gorm:"column:id"`
	CreateAt time.Time `gorm:"column:created_at"`
	UpdateAt time.Time `gorm:"column:updated_at"`
}
