package cves

import (
	"app/base/logging"
	"app/base/utils"
	"app/manager/base"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GetExposedClustersSelect
// @Description CVE exposed clusters data
// @Description presents in response
type GetExposedClustersSelect struct {
	UUID     string `json:"uuid"`
	Status   string `json:"status"`
	Version  string `json:"version"`
	Provider string `json:"provider"`
}

var (
	getExposedClustersAllowedFilters = []string{
		base.SortQuery,
		base.LimitQuery,
		base.OffsetQuery,
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

func init() {
	var err error
	logger, err = logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// GetExposedClusters represents exposed clusters endpoint controller.
//
// @id GetExposedClusters
// @summary List of exposed clusters for CVE
// @Tags cves
// @description Endpoint return exposed clusters for given CVE
// @accept */*
// @produce json
// @Param cve_name path  string true  "CVE name"
// @Param sort     query []string false "column for sort"      collectionFormat(multi) collectionFormat(csv)
// @Param search   query string   false "cve name/desc search" example(CVE-2021-)
// @Param limit    query int      false "limit per page"       example(10)
// @Param offset   query int      false "page offset"          example(10)
// @router /cves/{cve_name}/exposed_clusters [get]
// @success 200 {object} base.Response{data=GetExposedClustersSelect}
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "{cve_name} not found"
// @failure 500 {object} base.Error
func (c *Controller) GetExposedClusters(ctx *gin.Context) {
	cveName := ctx.Param("cve_name")
	accountID := ctx.GetInt64("account_id")

	filters := base.GetRequestedFilters(ctx)

	query := c.BuildExposedClustersQuery(cveName, accountID)
	err := base.ApplyFilters(query, getExposedClustersAllowedFilters, filters, getExposedClustersFilterArgs)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	var exposedClusters GetExposedClustersSelect
	result := query.First(&exposedClusters)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(
				http.StatusNotFound,
				base.BuildErrorResponse(http.StatusNotFound, fmt.Sprintf("%s not found", cveName)),
			)
			return
		}
		logger.Errorf("Database error: %s", result.Error)
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		return
	}

	ctx.JSON(http.StatusOK, base.BuildResponse(exposedClusters, base.BuildMeta(make(map[string]base.Filter), getExposedClustersAllowedFilters)))
}

func (c *Controller) BuildExposedClustersQuery(cveName string, accountID int64) *gorm.DB {
	return c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.status, cluster.version, cluster.provider`).
		Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON image_cve.cve_id = cve.id").
		Where("cve.name = ?", cveName).
		Where("cluster.account_id = ?", accountID)
}
