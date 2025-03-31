package entity

type WorkspaceEntity struct {
	BaseEntity
	Name                string
	Description         string
	Code                string
	WorkspaceUserEntity []WorkspaceUserEntity
}
