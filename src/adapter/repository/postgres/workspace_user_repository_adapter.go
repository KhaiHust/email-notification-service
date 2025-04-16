package postgres

import (
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type WorkspaceUserRepositoryAdapter struct {
	base
}

func NewWorkspaceUserRepositoryAdapter(db *gorm.DB) port.IWorkspaceUserRepositoryPort {
	return &WorkspaceUserRepositoryAdapter{
		base: base{db},
	}
}
