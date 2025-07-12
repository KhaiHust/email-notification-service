package apihelper

import "fmt"

type PageCursor struct {
	Since     int64  `json:"since"`
	Until     int64  `json:"until"`
	SortOrder string `json:"sort_order"` // "ASC" or "DESC"
	RawQuery  string `json:"raw_query"`  // e.g., "since=123&limit=20"
}

type PagingMetadata struct {
	Since        *int64      `json:"since"`
	Until        *int64      `json:"until"`
	Limit        int64       `json:"limit"`
	TotalItems   int64       `json:"total_items"`
	PageSize     int         `json:"page_size"`
	TotalPages   int64       `json:"total_pages"`
	HasMore      bool        `json:"has_more"`
	NextPage     *PageCursor `json:"next_page,omitempty"`
	PreviousPage *PageCursor `json:"previous_page,omitempty"`
}

func BuildIDPaginatedResponse[T any](
	items []T,
	since, until, limit, totalItems *int64,
	getID func(T) int64,
	sortOrder string, // "ASC" or "DESC"
) PagingMetadata {
	pageSize := len(items)
	hasMore := int64(pageSize) == *limit

	var nextPage, prevPage *PageCursor
	if pageSize > 0 {
		firstID := getID(items[0])
		lastID := getID(items[pageSize-1])

		switch sortOrder {
		case "ASC":
			nextPage = &PageCursor{
				Since:     lastID,
				Until:     0,
				SortOrder: "ASC",
				RawQuery:  fmt.Sprintf("since=%d&limit=%d&sort_order=ASC", lastID, limit),
			}
			prevPage = &PageCursor{
				Since:     0,
				Until:     firstID,
				SortOrder: "ASC",
				RawQuery:  fmt.Sprintf("until=%d&limit=%d&sort_order=ASC", firstID, limit),
			}
		case "DESC":
			nextPage = &PageCursor{
				Since:     0,
				Until:     lastID,
				SortOrder: "DESC",
				RawQuery:  fmt.Sprintf("until=%d&limit=%d&sort_order=DESC", lastID, limit),
			}
			prevPage = &PageCursor{
				Since:     firstID,
				Until:     0,
				SortOrder: "DESC",
				RawQuery:  fmt.Sprintf("since=%d&limit=%d&sort_order=DESC", firstID, limit),
			}
		}
	}

	return PagingMetadata{
		Since:        since,
		Until:        until,
		Limit:        *limit,
		TotalItems:   *totalItems,
		PageSize:     pageSize,
		TotalPages:   (*totalItems + *limit - 1) / *limit,
		HasMore:      hasMore,
		NextPage:     nextPage,
		PreviousPage: prevPage,
	}
}
