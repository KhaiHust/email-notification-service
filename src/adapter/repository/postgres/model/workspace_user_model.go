package model

type WorkspaceUserModel struct {
	BaseModel
	WorkspaceID int64  `gorm:"column:workspace_id"`
	UserID      int64  `gorm:"column:user_id"`
	Role        string `gorm:"column:role"`
}

func (WorkspaceUserModel) TableName() string {
	return "workspace_users"
}
