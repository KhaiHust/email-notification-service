package entity

type ApiKeyEntity struct {
	BaseEntity
	WorkspaceID int64
	Name        string
	KeyHash     string
	RawPrefix   string
	Environment string
	ExpiresAt   *int64
	Revoked     bool
}
