package middlewares

import (
	"app/manager/base"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseCommaParams Golang HTTP module does not parse query
// array arguments such as -> filter=1,2,3&filter=4 -> ["1,2,3", "4"]
// which should be ["1", "2", "3", "4"]
func ParseCommaParams(values []string) []string {
	var res []string
	for _, val := range values {
		vals := strings.Split(val, ",")
		res = append(res, vals...)
	}
	return res
}

func ApplyDefaultFilters(requestedFilters map[string]base.Filter) {
	// set default limit paging if not set
	if _, exists := requestedFilters[base.LimitQuery]; !exists {
		requestedFilters[base.LimitQuery] = &base.Limit{
			RawFilter: base.RawFilter{RawParam: "limit", RawValues: []string{"20"}},
			Value:     uint64(20)}
	}
	// set default offset paging if not set
	if _, exists := requestedFilters[base.OffsetQuery]; !exists {
		requestedFilters[base.OffsetQuery] = &base.Offset{
			RawFilter: base.RawFilter{RawParam: "offset", RawValues: []string{"0"}},
			Value:     uint64(0)}
	}
	// always sort, at least default values
	if _, exists := requestedFilters[base.SortQuery]; !exists {
		requestedFilters[base.SortQuery] = &base.Sort{
			RawFilter: base.RawFilter{RawParam: "sort", RawValues: []string{}},
			Values:    []base.SortItem{},
		}
	}
}

func Filterer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filters := make(map[string]base.Filter)
		for param, rawValues := range ctx.Request.URL.Query() {
			values := ParseCommaParams(rawValues)
			filter, err := base.ParseFilter(param, values)
			if err != nil && errors.Is(err, base.ErrInvalidFilterArgument) {
				continue
			} else if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
				return
			}
			filters[filter.RawQueryName()] = filter
		}
		ApplyDefaultFilters(filters)
		ctx.Set("filters", filters)
	}
}
