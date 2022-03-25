package base

type Error struct {
	Error ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

func BuildErrorResponse(status int, detail string) Error {
	return Error{Error: ErrorDetail{Detail: detail, Status: status}}
}
