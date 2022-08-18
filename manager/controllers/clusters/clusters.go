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
	CriticalCount  *int64 `json:"critical"`
	ImportantCount *int64 `json:"important"`
	ModerateCount  *int64 `json:"moderate"`
	LowCount       *int64 `json:"low"`
}

// GetClustersSelect
// @Description clusters data
type GetClustersSelect struct {
	UUID        string                `json:"id"`
	DisplayName string                `json:"display_name"`
	Status      string                `json:"status"`
	Type        string                `json:"type"`
	Version     string                `json:"version"`
	Provider    string                `json:"provider"`
	Region      string                `json:"region"`
	Severities  *ClusterCveSeverities `json:"cves_severity" gorm:"embedded"`
	LastSeen    time.Time             `json:"last_seen"`
}

type GetClustersResponse struct {
	Data []GetClustersSelect `json:"data"`
	Meta interface{}         `json:"meta"`
}

var (
	getClustersAllowedFilters = []string{
		base.SearchQuery,
		base.ClusterSeverityQuery,
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
				"display_name": "display_name"},
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
// @Param limit    			query int      false "limit per page"           example(10) minimum(0) maximum(100)
// @Param offset   			query int      false "page offset"              example(10) minimum(0)
// @Param cluster_severity  query []string false "array of severity names"  enums(Low,Moderate,Important,Critical)
// @router /clusters [get]
// @success 200 {object} GetClustersResponse
// @failure 400 {object} base.Error
// @failure 500 {object} base.Error
func (c *Controller) GetClusters(ctx *gin.Context) {
	var clusterIDs []string
	var clusterInfoMap map[string]amsclient.ClusterInfo
	var err error
	if utils.Cfg.AmsEnabled {
		orgID := ctx.GetString("org_id")
		clusterInfoMap, err = c.AMSClient.GetClustersForOrganization(orgID, nil, nil)
		if err != nil {
			c.Logger.Errorf("Error returned from AMS client: %s", err.Error())
			ctx.AbortWithStatusJSON(http.StatusBadGateway, base.BuildErrorResponse(http.StatusBadGateway, "Error returned from AMS API"))
			return
		}
		for clusterID := range clusterInfoMap {
			clusterIDs = append(clusterIDs, clusterID)
		}
	}

	accountID := ctx.GetInt64("account_id")
	filters := base.GetRequestedFilters(ctx)

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

	// Set cluster details fetched from AMS API
	if utils.Cfg.AmsEnabled {
		fullClustersData := []GetClustersSelect{}
		for _, clusterRow := range clustersData {
			if clusterInfo, ok := clusterInfoMap[clusterRow.UUID]; ok {
				clusterRow.DisplayName = clusterInfo.DisplayName
				clusterRow.Status = base.EmptyToNA(clusterInfo.Status)
				clusterRow.Type = base.EmptyToNA(clusterInfo.Type)
				clusterRow.Version = base.EmptyToNA(clusterInfo.Version)
				clusterRow.Provider = base.EmptyToNA(clusterInfo.Provider)
				clusterRow.Region = base.EmptyToNA(clusterInfo.Region)
			}
			fullClustersData = append(fullClustersData, clusterRow)
		}
		clustersData = fullClustersData
	}

	ctx.JSON(http.StatusOK, GetClustersResponse{clustersData, base.BuildMeta(usedFilters, &totalItems)})
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
		Select(`cluster.uuid, cluster.uuid AS display_name, cluster.status, cluster.version, cluster.provider,
				COALESCE(cc, 0) AS critical_count, COALESCE(ic, 0) AS important_count,
				COALESCE(mc, 0) AS moderate_count, COALESCE(lc, 0) AS low_count,
				cluster.last_seen`).
		Joins("JOIN (?) AS cluster_cves ON cluster.id = cluster_cves.id", subquery).
		Where("cluster.account_id = ?", accountID)
}
