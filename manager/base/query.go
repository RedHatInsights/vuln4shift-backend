package base

import (
	"gorm.io/gorm"
)

// Performs Find on query to given response, and overrides the limit and offsets
// to return the total count of items
func ListQueryFind(tx *gorm.DB, response interface{}) (*gorm.DB, int64) {
	res := tx.Find(response)
	if res.Error != nil {
		return res, 0
	}
	var count int64
	tx.Limit(-1).Offset(-1).Count(&count)
	return res, count
}
