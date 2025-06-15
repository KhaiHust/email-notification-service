package model

type WorkspaceModel struct {
	BaseModel
	Name               string                `gorm:"column:name"`
	Description        string                `gorm:"column:description"`
	Code               string                `gorm:"column:code"`
	WorkspaceUserModel []*WorkspaceUserModel `gorm:"foreignKey:WorkspaceID"`
}

func (WorkspaceModel) TableName() string {
	return "workspaces"
}
