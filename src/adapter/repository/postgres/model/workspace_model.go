package model

type WorkspaceModel struct {
	BaseModel
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (WorkspaceModel) TableName() string {
	return "workspaces"
}
