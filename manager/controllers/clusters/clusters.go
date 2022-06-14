package clusters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"app/base/models"
	"app/manager/base"
)

type ClusterCveSeverities struct {
	CriticalCount  *int64 `json:"critical"`
	ImportantCount *int64 `json:"important"`
	ModerateCount  *int64 `json:"moderate"`
	LowCount       *int64 `json:"low"`
}

// GetClustersSelect
// @Description clusters data
type GetClustersSelect struct {
	UUID       *string               `json:"id"`
	Status     *string               `json:"status"`
	Version    *string               `json:"version"`
	Provider   *string               `json:"provider"`
	Severities *ClusterCveSeverities `json:"cves_severity" gorm:"embedded"`
}

type GetClustersResponse []GetClustersSelect

var (
	getClustersAllowedFilters = []string{
		base.SortQuery,
		base.LimitQuery,
		base.OffsetQuery,
		base.SearchQuery,
		base.ClusterSeverityQuery,
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
// @Param sort     			query []string false "column for sort"          collectionFormat(multi) collectionFormat(csv)
// @Param search   			query string   false "cluster UUID search"      example(123e4567-e89b-12d3-a456-426614174000)
// @Param limit    			query int      false "limit per page"           example(10)
// @Param offset   			query int      false "page offset"              example(10)
// @Param cluster_severity  query []string false "array of severity names"  enums(Low,Moderate,Important,Critical)
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
	result, totalItems := base.ListQueryFind(query, &clustersData)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", result.Error)
		return
	}

	ctx.JSON(http.StatusOK, base.BuildResponse(clustersData, base.BuildMeta(filters, getClustersAllowedFilters, &totalItems)))
}

func (c *Controller) BuildClustersQuery(accountID int64) *gorm.DB {
	subquery := c.Conn.Table("cluster").
		Select(`cluster.id,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS critical_count,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS important_count,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS moderate_count,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS low_count`,
			models.Critical, models.Important, models.Moderate, models.Low).
		Joins("LEFT JOIN cluster_image ON (cluster.id = cluster_image.cluster_id AND cluster.account_id = ?)", accountID).
		Joins("LEFT JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("LEFT JOIN cve ON image_cve.cve_id = cve.id").
		Group("cluster.id")

	return c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.status, cluster.version, cluster.provider,
				critical_count, important_count, moderate_count, low_count`).
		Joins("LEFT JOIN (?) AS cluster_cves ON cluster.id = cluster_cves.id", subquery)
}
