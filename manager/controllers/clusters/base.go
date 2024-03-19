package clusters

import (
	"app/manager/base"

	"github.com/google/uuid"
)

type Controller struct {
	base.Controller
}

// ClusterExists, checks if cluster exists in db with given accid and clusterid
func (c *Controller) ClusterExists(accountID int64, clusterID uuid.UUID) (bool, error) {
	res := c.Conn.Table("cluster").Where("account_id = ? AND uuid = ?", accountID, clusterID).Limit(1).Find(&struct{}{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
