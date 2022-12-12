package cves

import (
	"app/manager/amsclient"
	"app/manager/base"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GetExposedClustersCountResponse struct {
	Count int64 `json:"count"`
}

// GetExposedClustersCount represents exposed clusters count endpoint controller.
//
// @id GetExposedClustersCount
// @summary Quantity of exposed clusters for a CVE
// @security RhIdentity || BasicAuth
// @Tags cves
// @description Endpoint returns the number of exposed clusters for a given CVE
// @accept */*
// @produce json
// @Param cve_name    path  string   true  "CVE name"
// @router /cves/{cve_name}/exposed_clusters_count [get]
// @success 200 {object} GetExposedClustersCountResponse
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "{cve_name} not found"
// @failure 500 {object} base.Error
func (c *Controller) GetExposedClustersCount(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")
	orgID := ctx.GetString("org_id")
	cveName := ctx.Param("cve_name")
	filters := base.GetRequestedFilters(ctx)

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

	_, _, _, _, err := amsclient.DBFetchClusterDetails(c.Conn, c.AMSClient, accountID, orgID, false, &cveName)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Error fetching AMS data: %s", err.Error())
		return
	}

	query = c.BuildExposedClustersCountQuery(cveName, accountID)

	_, totalItems, inputErr, dbErr := base.ListQuery(query, nil, filters, nil, nil)
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

	ctx.JSON(http.StatusOK, GetExposedClustersCountResponse{Count: totalItems})
}

func (c *Controller) BuildExposedClustersCountQuery(cveName string, accountID int64) *gorm.DB {
	query := c.Conn.Table("cluster").
		Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON image_cve.cve_id = cve.id").
		Group("cluster.id").
		Where("cve.name = ?", cveName).
		Where("cluster.account_id = ?", accountID)

	return query
}
