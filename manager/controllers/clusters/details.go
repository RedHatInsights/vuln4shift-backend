package clusters

import (
	"app/manager/base"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetClusterDetailsSelect
// @Description Cluster details
// @Description presents in response
type GetClusterDetailsSelect struct {
	UUID        string    `json:"id"`
	DisplayName string    `json:"display_name"`
	LastSeen    time.Time `json:"last_seen"`
}

type GetClusterDetailsResponse struct {
	Data GetClusterDetailsSelect `json:"data"`
	Meta interface{}             `json:"meta"`
}

// GetClusterDetails represents /clusters/{cluster_id} endpoint controller.
//
// @id GetClusterDetails
// @summary Cluster details
// @security RhIdentity || BasicAuth
// @Tags clusters
// @description Endpoint returning details of the given single cluster
// @accept */*
// @produce json
// @Param cluster_id      path  string   true  "cluster ID"
// @router /clusters/{cluster_id} [get]
// @success 200 {object} GetClusterDetailsResponse
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "cluster does not exist"
// @failure 500 {object} base.Error
func (c *Controller) GetClusterDetails(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")
	clusterID, err := base.GetParamUUID(ctx, "cluster_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	query := c.BuildClusterDetailsQuery(accountID, clusterID)

	var clusterDetails GetClusterDetailsSelect
	result := query.First(&clusterDetails)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				base.BuildErrorResponse(http.StatusNotFound, "cluster does not exist"),
			)
			return
		}
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", result.Error)
		return
	}

	ctx.JSON(http.StatusOK,
		GetClusterDetailsResponse{clusterDetails, base.BuildMeta(make(map[string]base.Filter), nil)},
	)
}

func (c *Controller) BuildClusterDetailsQuery(accountID int64, clusterID uuid.UUID) *gorm.DB {
	return c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.uuid AS display_name, cluster.last_seen`).
		Where("cluster.account_id = ? AND cluster.uuid = ?", accountID, clusterID)
}
