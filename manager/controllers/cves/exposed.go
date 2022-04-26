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

	logger *logrus.Logger
)

func init() {
	logLevel := utils.GetEnv("LOGGING_LEVEL", "INFO")
	var err error
	logger, err = logging.CreateLogger(logLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func (c *Controller) GetExposedClusters(ctx *gin.Context) {
	cveName := ctx.Param("cveName")

	filters := base.GetRequestedFilters(ctx)

	query := c.BuildExposedClustersQuery(cveName)
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

func (c *Controller) BuildExposedClustersQuery(cveName string) *gorm.DB {
	return c.Conn.Table("cluster").
		Select(`cluster.uuid, cluster.status, cluster.version, cluster.provider`).
		Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON image_cve.cve_id = cve.id").
		Where("cve.name = ?", cveName)
}
