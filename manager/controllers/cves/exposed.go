package cves

import (
	"app/manager/base"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetExposedClustersSelect
// @Description CVE exposed clusters data
// @Description presents in response
type GetExposedClustersSelect struct {
	UUID        string `json:"id"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	Version     string `json:"version"`
	Provider    string `json:"provider"`
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
// @Param limit    query int      false "limit per page"       example(10)
// @Param offset   query int      false "page offset"          example(10)
// @router /cves/{cve_name}/exposed_clusters [get]
// @success 200 {object} GetExposedClustersResponse
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "{cve_name} not found"
// @failure 500 {object} base.Error
func (c *Controller) GetExposedClusters(ctx *gin.Context) {
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

	query = c.BuildExposedClustersQuery(cveName, accountID)

	exposedClusters := []GetExposedClustersSelect{}
	totalItems, inputErr, dbErr := base.ListQuery(query, getExposedClustersAllowedFilters, filters, getExposedClustersFilterArgs, &exposedClusters)
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

	ctx.JSON(http.StatusOK, GetExposedClustersResponse{exposedClusters, base.BuildMeta(make(map[string]base.Filter), getExposedClustersAllowedFilters, &totalItems)})
}

func (c *Controller) BuildExposedClustersQuery(cveName string, accountID int64) *gorm.DB {
	// FIXME: display_name is hardcoded to uuid
	return c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.uuid AS display_name, cluster.status, cluster.version, cluster.provider`).
		Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON image_cve.cve_id = cve.id").
		Where("cve.name = ?", cveName).
		Where("cluster.account_id = ?", accountID)
}
