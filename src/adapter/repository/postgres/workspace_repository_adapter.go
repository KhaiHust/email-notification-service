package postgres

import (
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type WorkspaceRepositoryAdapter struct {
	base
}

func NewWorkspaceRepositoryAdapter(db *gorm.DB) port.IWorkspaceRepositoryPort {
	return &WorkspaceRepositoryAdapter{
		base: base{db},
	}
}
