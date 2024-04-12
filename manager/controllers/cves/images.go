package cves

import (
	"app/base/utils"
	"app/manager/base"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var getCveImagesAllowedFilters = []string{
	base.SearchQuery,
	base.DataFormatQuery,
}

var getCveImagesFilterArgs = map[string]interface{}{
	base.SortFilterArgs: base.SortArgs{
		SortableColumns: map[string]string{
			"id":       "repository.id",
			"name":     "repository.repository",
			"registry": "repository.registry",
		},
		DefaultSortable: []base.SortItem{{Column: "id", Desc: false}},
	},
	base.SearchQuery: base.ImagesSearch,
}

// GetCveImagesSelect
// @Description Exposed images in cve data
// @Description presents in response
type GetCveImagesSelect struct {
	ExposedClusters *int32              `json:"exposed_clusters" csv:"exposed_clusters" gorm:"column:exposed_clusters"`
	Repository      *string             `json:"name" csv:"name" gorm:"column:repository"`
	Registry        *string             `json:"registry" csv:"registry" gorm:"column:registry"`
	Version         *utils.ImageVersion `json:"version" csv:"version" gorm:"column:tags"`
}

type GetCveImagesResponse struct {
	Data []GetCveImagesSelect `json:"data"`
	Meta interface{}          `json:"meta"`
}

// GetCveImages represents /cves/{cve_name}/exposed_images endpoint controller.
//
// @id GetCveImages
// @summary List of images affecting single cve
// @security RhIdentity || BasicAuth
// @Tags cves
// @description Endpoint returning images affecting the given single cluster
// @accept */*
// @produce json
// @Param cve_name        path  string   true  "cve ID"
// @Param sort            query []string false "column for sort"                                      collectionFormat(multi) collectionFormat(csv)
// @Param search          query string   false "image name/registry search"                           example(ubi8)
// @Param limit           query int      false "limit per page"                                       example(10) minimum(0) maximum(100)
// @Param offset          query int      false "page offset"                                          example(10) minimum(0)
// @Param data_format     query string   false "data section format"                                  enums(json,csv)
// @Param report          query bool     false "overrides limit and offset to return everything"
// @router /cves/{cve_name}/exposed_images [get]
// @success 200 {object} GetCveImagesResponse
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "cve does not exist"
// @failure 500 {object} base.Error
func (c *Controller) GetCveImages(ctx *gin.Context) {
	accountID := ctx.GetInt64("account_id")
	cveName := ctx.Param("cve_name")

	exists, err := c.CveExists(cveName)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			base.BuildErrorResponse(http.StatusInternalServerError, "Internal server error"),
		)
		c.Logger.Errorf("Database error: %s", err.Error())
		return
	} else if !exists {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			base.BuildErrorResponse(http.StatusNotFound, "Cve does not exist"),
		)
		return
	}

	filters := base.GetRequestedFilters(ctx)
	query := c.BuildCveImagesQuery(accountID, cveName)

	dataRes := []GetCveImagesSelect{}
	usedFilters, totalItems, inputErr, dbErr := base.ListQuery(query, getCveImagesAllowedFilters, filters, getCveImagesFilterArgs, &dataRes)
	if inputErr != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			base.BuildErrorResponse(http.StatusBadRequest, inputErr.Error()),
		)
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

	resp, err := base.BuildDataMetaResponse(dataRes, base.BuildMeta(usedFilters, &totalItems, nil, nil, nil), usedFilters)
	if err != nil {
		c.Logger.Errorf("Internal server error: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) BuildCveImagesQuery(accountID int64, cveName string) *gorm.DB {
	cntSubquery := c.Conn.Table("cluster").
		Select(`COUNT(DISTINCT cluster_image.cluster_id) AS ec, image_cve.image_id as image_id`).
		Joins("JOIN cluster_image ON cluster_image.cluster_id = cluster.id").
		Joins("JOIN repository_image ON repository_image.image_id = cluster_image.image_id").
		Joins("JOIN repository ON repository.id=repository_image.repository_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON cve.id=image_cve.cve_id").
		Group("image_cve.image_id").
		Where("cluster.account_id = ? AND cve.name = ?", accountID, cveName)

	return c.Conn.Table("repository").
		Select(`COALESCE(ec, 0) AS exposed_clusters, repository.repository, repository.registry, COALESCE(repository_image.tags, '[]') AS tags`).
		Joins("JOIN repository_image ON repository.id = repository_image.repository_id").
		Joins("JOIN cluster_image ON cluster_image.image_id=repository_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve ON cve.id=image_cve.cve_id").
		Joins("INNER JOIN (?) as cnt_subquery ON image_cve.image_id=cnt_subquery.image_id", cntSubquery).
		Where("cluster.account_id = ? AND cve.name = ?", accountID, cveName)
}
