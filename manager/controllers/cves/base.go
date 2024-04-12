package cves

import (
	"app/manager/base"
)

type Controller struct {
	base.Controller
}

func (c *Controller) CveExists(cveName string) (bool, error) {
	// Check if CVE exists first
	res := c.Conn.Table("cve").Where("name = ?", cveName).Limit(1).Find(&struct{}{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
