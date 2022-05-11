package cves

import (
	"app/base/models"
	"app/manager/base"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var getCvesAllowedFilters = []string{base.SortQuery, base.LimitQuery, base.OffsetQuery,
	base.SearchQuery, base.PublishedQuery, base.SeverityQuery, base.CvssScoreQuery,
	base.AffectedClustersQuery, base.AffectedImagesQuery}

var getCvesFilterArgs = map[string]interface{}{
	base.SortFilterArgs: base.SortArgs{
		SortableColumns: map[string]string{
			"id":               "cve.id",
			"cvss_score":       "COALESCE(cve.cvss3_score, cve.cvss2_score, 0.0)",
			"severity":         "cve.severity",
			"publish_date":     "cve.public_date",
			"synopsis":         "cve.name",
			"clusters_exposed": "clusters_exposed",
			"images_exposed":   "images_exposed",
		},
		DefaultSortable: []base.SortItem{{Column: "id", Desc: false}},
	},
	base.SearchQuery: base.CveSearch,
}

// GetCvesSelect
// @Description CVE in workload data
// @Description presents in response
type GetCvesSelect struct {
	Cvss2Score      *float32         `json:"cvss2_score"`
	Cvss3Score      *float32         `json:"cvss3_score"`
	Description     *string          `json:"description"`
	Severity        *models.Severity `json:"severity"`
	PublicDate      *time.Time       `json:"publish_date"`
	Name            *string          `json:"synopsis"`
	ClustersExposed *int64           `json:"clusters_exposed"`
	ImagesExposed   *int64           `json:"images_exposed"`
}

type GetCvesResponse []GetCvesSelect

// GetCves represents CVEs endpoint controller.
//
// @id GetCves
// @summary List of CVEs affecting the workload
// @security RhIdentity
// @Tags cves
// @description Endpoint returning CVEs affecting the current workload
// @accept */*
// @produce json
// @Param sort     query []string false "column for sort"      collectionFormat(multi) collectionFormat(csv)
// @Param search   query string   false "cve name/desc search" example(CVE-2021-)
// @Param limit    query int      false "limit per page"       example(10)
// @Param offset   query int      false "page offset"          example(10)
// @router /cves [get]
// @success 200 {object} base.Response{data=GetCvesResponse}
// @failure 400 {object} base.Error
func (c *Controller) GetCves(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")

	filters := base.GetRequestedFilters(ctx)
	query := c.BuildCvesQuery(accountID)

	err := base.ApplyFilters(query, getCvesAllowedFilters, filters, getCvesFilterArgs)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	dataRes := []GetCvesSelect{}
	query.Find(&dataRes)

	ctx.JSON(http.StatusOK, base.BuildResponse(dataRes, base.BuildMeta(filters, getCvesAllowedFilters)))
}

func (c *Controller) BuildCvesQuery(accountID int64) *gorm.DB {
	return c.Conn.Table("cve").
		Select(`cve.name, cve.description, cve.public_date, cve.severity,
							cve.cvss2_score, cve.cvss3_score,
							14 AS clusters_exposed, 8 AS images_exposed`).
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Where("cluster.account_id = ?", accountID).
		Group("cve.id, cve.name, cve.description, cve.public_date, cve.severity, cve.cvss3_score, cve.cvss2_score")
}
