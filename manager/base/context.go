package base

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetRequestedFilters gets requested parsed filters from gin context
// returns empty map if not exists
func GetRequestedFilters(ctx *gin.Context) map[string]Filter {
	if f, exists := ctx.Get("filters"); exists {
		if f, ok := f.(map[string]Filter); ok {
			return f
		}
	}
	return map[string]Filter{}
}

// GetUUID gets raw value from gin context and returns an parsed UUID (or error)
func GetParamUUID(ctx *gin.Context, key string) (uuid.UUID, error) {
	rawVal := ctx.Param(key)
	return uuid.Parse(rawVal)
}
