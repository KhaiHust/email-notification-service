package request

type GetApiKeyRequestFilter struct {
	WorkspaceIDs []int64
	Environments []string
}
