package clusters

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
		base.ProviderQuery,
		base.StatusQuery,
		base.VersionQuery,
	}

	getClustersFilterArgs = map[string]interface{}{
		base.SortFilterArgs: base.SortArgs{
			SortableColumns: map[string]string{
				"id":               "cluster.id",
				"status":           "cluster.status",
				"version":          "cluster.version",
				"provider":         "cluster.provider",
				"uuid":             "cluster.uuid",
				"last_seen":        "cluster.last_seen",
				"display_name":     "cluster.display_name",
				"type":             "cluster.type",
				"cluster_severity": "(cluster.cve_cache_critical,cluster.cve_cache_important,cluster.cve_cache_moderate,cluster.cve_cache_low)"},
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
// @Param sort     			query []string false "column for sort"                                      collectionFormat(multi) collectionFormat(csv)
// @Param search   			query string   false "cluster search"                                       example(123e4567-e89b-12d3-a456-426614174000)
// @Param limit    			query int      false "limit per page"                                       example(10) minimum(0) maximum(100)
// @Param offset   			query int      false "page offset"                                          example(10) minimum(0)
// @Param data_format       query string   false "data section format"                                  enums(json,csv)
// @Param report            query bool     false "overrides limit and offset to return everything"
// @Param provider          query []string false "provider of the cluster"
// @Param status            query []string false "status of the cluster"
// @Param version           query []string false "version of the cluster"
// @Param cluster_severity  query []string false "array of severity names"                              enums(Low,Moderate,Important,Critical)
// @router /clusters [get]
// @success 200 {object} GetClustersResponse
// @failure 400 {object} base.Error
// @failure 500 {object} base.Error
func (c *Controller) GetClusters(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")
	orgID := ctx.GetString("org_id")
	filters := base.GetRequestedFilters(ctx)

	clusterIDs, clusterStatuses, clusterVersions, clusterProviders, err := amsclient.DBFetchClusterDetails(c.Conn, c.AMSClient, accountID, orgID, utils.Cfg.AmsEnabled, nil)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Error fetching AMS data: %s", err.Error())
		return
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

	resp, err := base.BuildDataMetaResponse(clustersData, base.BuildMeta(usedFilters, &totalItems, &clusterStatuses, &clusterVersions, &clusterProviders, nil), usedFilters)
	if err != nil {
		c.Logger.Errorf("Internal server error: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) BuildClustersQuery(accountID int64, clusterIDs []string) *gorm.DB {
	query := c.Conn.Table("cluster").
		Select(`cluster.uuid,
				COALESCE(cluster.display_name, cluster.uuid::text) as display_name,
				COALESCE(cluster.status, 'N/A') as status,
				COALESCE(cluster.type, 'N/A') as type,
				COALESCE(cluster.version, 'N/A') as version,
				COALESCE(cluster.provider, 'N/A') as provider,
				cluster.cve_cache_critical AS critical_count,
				cluster.cve_cache_important AS important_count,
				cluster.cve_cache_moderate AS moderate_count,
				cluster.cve_cache_low AS low_count,
				cluster.last_seen`).
		Where("cluster.account_id = ?", accountID)

	if utils.Cfg.AmsEnabled {
		query = query.Where("cluster.uuid IN ?", clusterIDs)
	}

	return query
}
