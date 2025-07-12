package request

type BaseFilter struct {
	CreatedAtFrom *int64
	CreatedAtTo   *int64
	UpdatedAtFrom *int64
	UpdatedAtTo   *int64
	Limit         *int64
	Since         *int64
	Until         *int64
	SortOrder     string
}
