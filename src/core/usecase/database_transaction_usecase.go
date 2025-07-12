package usecase

import (
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type IDatabaseTransactionUseCase interface {
	StartTx() *gorm.DB
	CommitTx(tx *gorm.DB) error
	RollbackTx(tx *gorm.DB) error
}
type DatabaseTransactionUseCase struct {
	databaseTransactionPort port.IDatabaseTransactionPort
}

func (d DatabaseTransactionUseCase) StartTx() *gorm.DB {
	return d.databaseTransactionPort.StartTx()
}

func (d DatabaseTransactionUseCase) CommitTx(tx *gorm.DB) error {
	return d.databaseTransactionPort.Commit(tx)
}

func (d DatabaseTransactionUseCase) RollbackTx(tx *gorm.DB) error {
	return d.databaseTransactionPort.Rollback(tx)
}

func NewDatabaseTransactionUseCase(databaseTransactionPort port.IDatabaseTransactionPort) IDatabaseTransactionUseCase {
	return &DatabaseTransactionUseCase{
		databaseTransactionPort: databaseTransactionPort,
	}
}
