package clusters

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"app/base/models"
	"app/base/utils"
	"app/manager/amsclient"
	"app/manager/base"
)

type ClusterCveSeverities struct {
	CriticalCount  *int64 `json:"critical" csv:"critical"`
	ImportantCount *int64 `json:"important" csv:"important"`
	ModerateCount  *int64 `json:"moderate" csv:"moderate"`
	LowCount       *int64 `json:"low" csv:"low"`
}

// GetClustersSelect
// @Description clusters data
type GetClustersSelect struct {
	UUID        string                `json:"id" csv:"id"`
	DisplayName string                `json:"display_name" csv:"display_name"`
	Status      string                `json:"status" csv:"status"`
	Type        string                `json:"type" csv:"type"`
	Version     string                `json:"version" csv:"version"`
	Provider    string                `json:"provider" csv:"provider"`
	Severities  *ClusterCveSeverities `json:"cves_severity" csv:"-" gorm:"embedded"`
	LastSeen    time.Time             `json:"last_seen" csv:"last_seen"`
}

type GetClustersResponse struct {
	Data []GetClustersSelect `json:"data"`
	Meta interface{}         `json:"meta"`
}

var (
	getClustersAllowedFilters = []string{
		base.SearchQuery,
		base.ClusterSeverityQuery,
		base.DataFormatQuery,
	}

	getClustersFilterArgs = map[string]interface{}{
		base.SortFilterArgs: base.SortArgs{
			SortableColumns: map[string]string{
				"id":           "cluster.id",
				"status":       "cluster.status",
				"version":      "cluster.version",
				"provider":     "cluster.provider",
				"uuid":         "cluster.uuid",
				"last_seen":    "cluster.last_seen",
				"display_name": "cluster.display_name",
				"type":         "cluster.type"},
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
// @Param search   			query string   false "cluster search"           example(123e4567-e89b-12d3-a456-426614174000)
// @Param limit    			query int      false "limit per page"           example(10) minimum(0) maximum(100)
// @Param offset   			query int      false "page offset"              example(10) minimum(0)
// @Param data_format       query string   false "data section format"      enums(json,csv)
// @Param cluster_severity  query []string false "array of severity names"  enums(Low,Moderate,Important,Critical)
// @router /clusters [get]
// @success 200 {object} GetClustersResponse
// @failure 400 {object} base.Error
// @failure 500 {object} base.Error
func (c *Controller) GetClusters(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")
	filters := base.GetRequestedFilters(ctx)

	var clusterIDs []string
	var clusterInfoMap map[string]amsclient.ClusterInfo

	// Meta section sets
	var clusterStatuses, clusterVersions, clusterProviders map[string]struct{}

	var err error
	if utils.Cfg.AmsEnabled {
		orgID := ctx.GetString("org_id")
		clusterInfoMap, err = c.AMSClient.GetClustersForOrganization(orgID)
		if err != nil {
			c.Logger.Errorf("Error returned from AMS client: %s", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, base.BuildErrorResponse(http.StatusBadGateway, "Error returned from AMS API"))
			return
		}
		clusterIDs, clusterStatuses, clusterVersions, clusterProviders, err = amsclient.DBSyncClusterDetails(c.Conn, accountID, clusterInfoMap)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
			)
			c.Logger.Errorf("Database error: %s", err.Error())
			return
		}
	}

	query := c.BuildClustersQuery(accountID, clusterIDs)

	clustersData := []GetClustersSelect{}
	usedFilters, totalItems, inputErr, dbErr := base.ListQuery(query, getClustersAllowedFilters, filters, getClustersFilterArgs, &clustersData)
	if inputErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, inputErr.Error()))
		return
	}
	if dbErr != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", dbErr.Error())
		return
	}

	resp, err := base.BuildDataMetaResponse(clustersData, base.BuildMeta(usedFilters, &totalItems, &clusterStatuses, &clusterVersions, &clusterProviders), usedFilters)
	if err != nil {
		c.Logger.Errorf("Internal server error: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) BuildClustersQuery(accountID int64, clusterIDs []string) *gorm.DB {
	subquery := c.Conn.Table("cluster").
		Select(`cluster.id,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS cc,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS ic,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS mc,
				COUNT(DISTINCT CASE WHEN cve.severity = ? THEN cve.id ELSE NULL END) AS lc`,
			models.Critical, models.Important, models.Moderate, models.Low).
		Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON image_cve.cve_id = cve.id").
		Where("cluster.account_id = ?", accountID).
		Group("cluster.id")

	if utils.Cfg.AmsEnabled {
		subquery = subquery.Where("cluster.uuid IN ?", clusterIDs)
	}

	return c.Conn.Table("cluster").
		Select(`cluster.uuid,
				COALESCE(cluster.display_name, cluster.uuid::text) as display_name,
				COALESCE(cluster.status, 'N/A') as status,
				COALESCE(cluster.type, 'N/A') as type,
				COALESCE(cluster.version, 'N/A') as version,
				COALESCE(cluster.provider, 'N/A') as provider,
				COALESCE(cc, 0) AS critical_count, COALESCE(ic, 0) AS important_count,
				COALESCE(mc, 0) AS moderate_count, COALESCE(lc, 0) AS low_count,
				cluster.last_seen`).
		Joins("JOIN (?) AS cluster_cves ON cluster.id = cluster_cves.id", subquery).
		Where("cluster.account_id = ?", accountID)
}
