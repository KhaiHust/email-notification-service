package port

import "gorm.io/gorm"

type IDatabaseTransactionPort interface {
	StartTx() *gorm.DB
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
}
