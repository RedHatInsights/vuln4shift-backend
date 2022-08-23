package base

import (
	"errors"

	"github.com/gocarina/gocsv"
)

type Error struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

type DataMetaResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func BuildErrorResponse(status int, detail string) Error {
	return Error{Error: ErrorDetail{Detail: detail, Status: status}}
}

func BuildDataMetaResponse(data interface{}, meta interface{}, filters map[string]Filter) (DataMetaResponse, error) {
	if filter, exists := filters[DataFormatQuery]; exists {
		dataFormat, ok := filter.(*DataFormat)
		if !ok {
			return DataMetaResponse{}, errors.New("Invalid data format filter")
		}

		switch dataFormat.Value {
		case CSVFormat:
			data, err := gocsv.MarshalString(data)
			if err != nil {
				return DataMetaResponse{}, err
			}
			return DataMetaResponse{data, meta}, nil
		}
	}
	return DataMetaResponse{data, meta}, nil
}

// BuildMeta creates Meta section in response from requested filters
// result is map with query args and their raw values
func BuildMeta(requestedFilters map[string]Filter, totalItems *int64) map[string]interface{} {
	meta := make(map[string]interface{})
	for _, filter := range requestedFilters {
		meta[filter.RawQueryName()] = filter.RawQueryVal()
	}
	if totalItems != nil {
		meta["total_items"] = *totalItems
	}
	return meta
}

func EmptyToNA(input string) string {
	if input == "" {
		return "N/A"
	}
	return input
}
