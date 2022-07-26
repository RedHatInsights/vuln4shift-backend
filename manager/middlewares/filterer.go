package middlewares

import (
	"app/manager/base"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
			filter, err := base.ParseFilter(param, rawValues)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
				return
			}
			filters[filter.RawQueryName()] = filter
		}
		ApplyDefaultFilters(filters)
		ctx.Set("filters", filters)
	}
}
