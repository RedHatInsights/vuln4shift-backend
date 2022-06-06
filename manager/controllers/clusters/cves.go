package clusters

import (
	"app/base/models"
	"app/manager/base"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var getClusterCvesAllowedFilters = []string{base.SortQuery, base.LimitQuery, base.OffsetQuery,
	base.SearchQuery, base.PublishedQuery, base.SeverityQuery, base.CvssScoreQuery}

var getClusterCvesFilterArgs = map[string]interface{}{
	base.SortFilterArgs: base.SortArgs{
		SortableColumns: map[string]string{
			"id":             "cve.id",
			"cvss_score":     "COALESCE(cve.cvss3_score, cve.cvss2_score, 0.0)",
			"severity":       "cve.severity",
			"publish_date":   "cve.public_date",
			"synopsis":       "cve.name",
			"images_exposed": "images_exposed",
		},
		DefaultSortable: []base.SortItem{{Column: "id", Desc: false}},
	},
	base.SearchQuery: base.CveSearch,
}

// GetClusterCvesSelect
// @Description CVE in cluster data
// @Description presents in response
type GetClusterCvesSelect struct {
	Cvss2Score    *float32         `json:"cvss2_score"`
	Cvss3Score    *float32         `json:"cvss3_score"`
	Description   *string          `json:"description"`
	Severity      *models.Severity `json:"severity"`
	PublicDate    *time.Time       `json:"publish_date"`
	Name          *string          `json:"synopsis"`
	ImagesExposed *int64           `json:"images_exposed"`
}

type GetClusterCvesResponse []GetClusterCvesSelect

// GetClusterCves represents /clusters/{cluster_id}/cves endpoint controller.
//
// @id GetClusterCves
// @summary List of CVEs affecting/nonaffecting single cluster
// @security RhIdentity || BasicAuth
// @Tags clusters
// @description Endpoint returning CVEs affecting/nonaffecting the given single cluster
// @accept */*
// @produce json
// @Param cluster_id      path  string   true  "cluster ID"
// @Param sort            query []string false "column for sort"                                      collectionFormat(multi) collectionFormat(csv)
// @Param search          query string   false "cve name/desc search"                                 example(CVE-2021-)
// @Param limit           query int      false "limit per page"                                       example(10)
// @Param offset          query int      false "page offset"                                          example(10)
// @Param published       query []string false "CVE publish date: (from date),(to date)"              collectionFormat(multi) collectionFormat(csv) minItems(2) maxItems(2)
// @Param severity        query []string false "array of severity names"                              enums(NotSet,None,Low,Medium,Moderate,Important,High,Critical)
// @Param cvss_score      query []number false "CVSS score of CVE: (from float),(to float)"           collectionFormat(multi) collectionFormat(csv) minItems(2) maxItems(2)
// @router /clusters/{cluster_id}/cves [get]
// @success 200 {object} base.Response{data=GetClusterCvesSelect}
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "cluster does not exist"
// @failure 500 {object} base.Error
func (c *Controller) GetClusterCves(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")
	clusterID, err := base.GetParamUUID(ctx, "cluster_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	exists, err := c.ClusterExists(accountID, clusterID)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", err.Error())
		return
	} else if !exists {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			base.BuildErrorResponse(http.StatusNotFound, "cluster does not exist"),
		)
		return
	}

	filters := base.GetRequestedFilters(ctx)
	query := c.BuildClusterCvesQuery(accountID, clusterID)

	err = base.ApplyFilters(query, getClusterCvesAllowedFilters, filters, getClusterCvesFilterArgs)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			base.BuildErrorResponse(http.StatusBadRequest, err.Error()),
		)
		return
	}
	dataRes := []GetClusterCvesSelect{}
	res, totalItems := base.ListQueryFind(query, &dataRes)
	if res.Error != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", res.Error)
		return
	}

	ctx.JSON(http.StatusOK,
		base.BuildResponse(dataRes, base.BuildMeta(filters, getClusterCvesAllowedFilters, &totalItems)),
	)
}

// ClusterExists, checks if cluster exists in db with given accid and clusterid
func (c *Controller) ClusterExists(accountID int64, clusterID uuid.UUID) (bool, error) {
	res := c.Conn.Table("cluster").Where("account_id = ? AND uuid = ?", accountID, clusterID).Limit(1).Find(&struct{}{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (c *Controller) BuildClusterCvesQuery(accountID int64, clusterID uuid.UUID) *gorm.DB {
	return c.Conn.Table("cve").
		Select(`cve.cvss2_score, cve.cvss3_score, cve.description, cve.severity,
			cve.public_date, cve.name,
			COUNT(DISTINCT cluster_image.image_id) as images_exposed`).
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Group("cve.id").
		Where("cluster.account_id = ? AND cluster.uuid = ?", accountID, clusterID)
}
