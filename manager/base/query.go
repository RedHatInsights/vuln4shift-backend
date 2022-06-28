package base

import (
	"gorm.io/gorm"
)

// ListQuery applies filters, queries and gets total count for listable endpoint
func ListQuery(tx *gorm.DB, allowedFilters []string, filters map[string]Filter, filterArgs map[string]interface{}, result interface{}) (totalItems int64, inputError error, dbError error) {
	inputError = ApplyFilters(tx, allowedFilters, filters, filterArgs)
	if inputError != nil {
		return
	}
	res := tx.Count(&totalItems)
	if res.Error != nil {
		dbError = res.Error
		return
	}
	spqFilters := []string{SortQuery, LimitQuery, OffsetQuery}
	inputError = ApplyFilters(tx, spqFilters, filters, filterArgs)
	if inputError != nil {
		return
	}
	res = tx.Find(result)
	if res.Error != nil {
		dbError = res.Error
		return
	}
	return
}
