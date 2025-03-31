package postgres

import (
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type DatabaseTransactionAdapter struct {
	base
}

func (d DatabaseTransactionAdapter) StartTx() *gorm.DB {
	return d.StartTransaction()
}

func (d DatabaseTransactionAdapter) Commit(tx *gorm.DB) error {
	return d.CommitTransaction(tx)
}

func (d DatabaseTransactionAdapter) Rollback(tx *gorm.DB) error {
	return d.RollbackTransaction(tx)
}

func NewDatabaseTransactionAdapter(db *gorm.DB) port.IDatabaseTransactionPort {
	return &DatabaseTransactionAdapter{
		base: base{db: db},
	}
}
