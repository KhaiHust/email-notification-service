package entity

type WorkspaceUserEntity struct {
	BaseEntity
	WorkspaceID int64
	UserID      int64
	Role        string
}
