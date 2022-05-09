package base

type Error struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

type Response struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func BuildErrorResponse(status int, detail string) Error {
	return Error{Error: ErrorDetail{Detail: detail, Status: status}}
}

// BuildMeta creates Meta section in response from requested filters
// result is map with query args and their raw values
func BuildMeta(requestedFilters map[string]Filter, allowedFilters []string) map[string]string {
	meta := make(map[string]string)
	for _, allowedFilter := range allowedFilters {
		if filter, requested := requestedFilters[allowedFilter]; requested {
			meta[filter.RawQueryName()] = filter.RawQueryVal()
		}
	}
	return meta
}

func BuildResponse(data interface{}, meta interface{}) Response {
	return Response{data, meta}
}
