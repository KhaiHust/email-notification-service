package request

type SendVolumeFilter struct {
	StartDate   *int64
	EndDate     *int64
	WorkspaceId int64
}
