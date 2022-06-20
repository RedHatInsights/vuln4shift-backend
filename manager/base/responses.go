package base

type Error struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

func BuildErrorResponse(status int, detail string) Error {
	return Error{Error: ErrorDetail{Detail: detail, Status: status}}
}

// BuildMeta creates Meta section in response from requested filters
// result is map with query args and their raw values
func BuildMeta(requestedFilters map[string]Filter, allowedFilters []string, totalItems *int64) map[string]interface{} {
	meta := make(map[string]interface{})
	for _, allowedFilter := range allowedFilters {
		if filter, requested := requestedFilters[allowedFilter]; requested {
			meta[filter.RawQueryName()] = filter.RawQueryVal()
		}
	}
	if totalItems != nil {
		meta["total_items"] = *totalItems
	}
	return meta
}
