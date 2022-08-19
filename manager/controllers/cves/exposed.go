package cves

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"app/base/utils"
	"app/manager/amsclient"
	"app/manager/base"
)

// GetExposedClustersSelect
// @Description CVE exposed clusters data
// @Description presents in response
type GetExposedClustersSelect struct {
	UUID        string     `json:"id"`
	DisplayName string     `json:"display_name"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
	Version     string     `json:"version"`
	Provider    string     `json:"provider"`
	Region      string     `json:"region"`
	LastSeen    *time.Time `json:"last_seen"`
}

type GetExposedClustersResponse struct {
	Data []GetExposedClustersSelect `json:"data"`
	Meta interface{}                `json:"meta"`
}

var (
	getExposedClustersAllowedFilters = []string{
		base.SearchQuery,
	}

	getExposedClustersFilterArgs = map[string]interface{}{
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

// GetExposedClusters represents exposed clusters endpoint controller.
//
// @id GetExposedClusters
// @summary List of exposed clusters for CVE
// @security RhIdentity || BasicAuth
// @Tags cves
// @description Endpoint return exposed clusters for given CVE
// @accept */*
// @produce json
// @Param cve_name path  string   true  "CVE name"
// @Param sort     query []string false "column for sort"      collectionFormat(multi) collectionFormat(csv)
// @Param search   query string   false "cluster search"       example(00000000-0000-0000-0000-000000000022)
// @Param limit    query int      false "limit per page"       example(10) minimum(0) maximum(100)
// @Param offset   query int      false "page offset"          example(10) minimum(0)
// @router /cves/{cve_name}/exposed_clusters [get]
// @success 200 {object} GetExposedClustersResponse
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "{cve_name} not found"
// @failure 500 {object} base.Error
func (c *Controller) GetExposedClusters(ctx *gin.Context) {
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

	cveName := ctx.Param("cve_name")
	accountID := ctx.GetInt64("account_id")

	// Check if CVE exists first
	query := c.BuildCveDetailsQuery(cveName)
	var cveDetails GetCveDetailsSelect
	result := query.First(&cveDetails)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				base.BuildErrorResponse(http.StatusNotFound, fmt.Sprintf("%s not found", cveName)),
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

	// If yes, select clusters
	filters := base.GetRequestedFilters(ctx)

	query = c.BuildExposedClustersQuery(cveName, accountID, clusterIDs)

	exposedClusters := []GetExposedClustersSelect{}
	_, totalItems, inputErr, dbErr := base.ListQuery(query, getExposedClustersAllowedFilters, filters, getExposedClustersFilterArgs, &exposedClusters)
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
		fullExposedClusters := []GetExposedClustersSelect{}
		for _, clusterRow := range exposedClusters {
			if clusterInfo, ok := clusterInfoMap[clusterRow.UUID]; ok {
				clusterRow.DisplayName = clusterInfo.DisplayName
				clusterRow.Status = base.EmptyToNA(clusterInfo.Status)
				clusterRow.Type = base.EmptyToNA(clusterInfo.Type)
				clusterRow.Version = base.EmptyToNA(clusterInfo.Version)
				clusterRow.Provider = base.EmptyToNA(clusterInfo.Provider)
				clusterRow.Region = base.EmptyToNA(clusterInfo.Region)
			}
			fullExposedClusters = append(fullExposedClusters, clusterRow)
		}
		exposedClusters = fullExposedClusters
	}

	ctx.JSON(http.StatusOK, GetExposedClustersResponse{exposedClusters, base.BuildMeta(make(map[string]base.Filter), &totalItems)})
}

func (c *Controller) BuildExposedClustersQuery(cveName string, accountID int64, clusterIDs []string) *gorm.DB {
	query := c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.uuid AS display_name, cluster.status, cluster.version, cluster.provider, cluster.last_seen,
		        COUNT(DISTINCT cluster_image.image_id) as images_exposed`).
		Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON image_cve.cve_id = cve.id").
		Group("cluster.id").
		Where("cve.name = ?", cveName).
		Where("cluster.account_id = ?", accountID)

	if utils.Cfg.AmsEnabled {
		query = query.Where("cluster.uuid IN ?", clusterIDs)
	}

	return query
}
