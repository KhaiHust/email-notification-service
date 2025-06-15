package request

type GetEmailProviderRequestFilter struct {
	Provider    *string
	WorkspaceID *int64
	Environment *string
}
