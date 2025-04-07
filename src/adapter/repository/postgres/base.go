package postgres

import (
	"gorm.io/gorm"
	"time"
)

type base struct {
	db *gorm.DB
}

func (b *base) StartTransaction() *gorm.DB {
	return b.db.Begin()
}
func (b *base) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}
func (b *base) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
func (b *base) BeforeUpdate() {
	b.db.Set("UpdatedAt", time.Now())
}
