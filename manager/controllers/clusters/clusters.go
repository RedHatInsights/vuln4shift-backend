package clusters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"app/manager/base"
)

// GetClustersSelect
// @Description clusters data
type GetClustersSelect struct {
	UUID     string `json:"uuid"`
	Status   string `json:"status"`
	Version  string `json:"version"`
	Provider string `json:"provider"`
}

type GetClustersResponse []GetClustersSelect

var (
	getClustersAllowedFilters = []string{
		base.SortQuery,
		base.LimitQuery,
		base.OffsetQuery,
		base.SearchQuery,
	}

	getClustersFilterArgs = map[string]interface{}{
		base.SortFilterArgs: base.SortArgs{
			SortableColumns: map[string]string{
				"id":       "cluster.id",
				"status":   "cluster.status",
				"version":  "cluster.version",
				"provider": "cluster.provider",
				"uuid":     "cluster.uuid"},
			DefaultSortable: []base.SortItem{{Column: "id", Desc: false}},
		},
		base.SearchQuery: base.ExposedClustersSearch,
	}
)

// GetClusters represents Clusters endpoint controller.
//
// @id GetClusters
// @summary List of Clusters
// @security RhIdentity || BasicAuth
// @Tags clusters
// @description Endpoint returning Clusters
// @accept */*
// @produce json
// @Param sort     			query []string false "column for sort"      collectionFormat(multi) collectionFormat(csv)
// @Param search   			query string   false "cluster UUID search"  example(123e4567-e89b-12d3-a456-426614174000)
// @Param limit    			query int      false "limit per page"       example(10)
// @Param offset   			query int      false "page offset"          example(10)
// @router /clusters [get]
// @success 200 {object} base.Response{data=GetClustersResponse}
// @failure 400 {object} base.Error
// @failure 500 {object} base.Error
func (c *Controller) GetClusters(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")

	filters := base.GetRequestedFilters(ctx)
	query := c.BuildClustersQuery(accountID)

	err := base.ApplyFilters(query, getClustersAllowedFilters, filters, getClustersFilterArgs)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	clustersData := []GetClustersSelect{}
	result := query.Find(&clustersData)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", result.Error)
		return
	}

	ctx.JSON(http.StatusOK, base.BuildResponse(clustersData, base.BuildMeta(filters, getClustersAllowedFilters)))

}

func (c *Controller) BuildClustersQuery(accountID int64) *gorm.DB {
	return c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.status, cluster.version, cluster.provider`).
		Where("cluster.account_id = ?", accountID)

}
