package clusters

import (
	"app/base/utils"
	"app/manager/base"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var getClusterImagesAllowedFilters = []string{
	base.SearchQuery,
	base.DataFormatQuery,
}

var getClusterImagesFilterArgs = map[string]interface{}{
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

// GetClusterImagesSelect
// @Description Exposed images in cluster data
// @Description presents in response
type GetClusterImagesSelect struct {
	Repository *string             `json:"name" csv:"name"`
	Registry   *string             `json:"registry" csv:"registry"`
	Version    *utils.ImageVersion `json:"version" csv:"version" gorm:"column:tags"`
}

type GetClusterImagesResponse struct {
	Data []GetClusterImagesSelect `json:"data"`
	Meta interface{}              `json:"meta"`
}

// GetClusterImages represents /clusters/{cluster_id}/exposed_images endpoint controller.
//
// @id GetClusterImages
// @summary List of images affecting single cluster
// @security RhIdentity || BasicAuth
// @Tags clusters
// @description Endpoint returning images affecting the given single cluster
// @accept */*
// @produce json
// @Param cluster_id      path  string   true  "cluster ID"
// @Param sort            query []string false "column for sort"                                      collectionFormat(multi) collectionFormat(csv)
// @Param search          query string   false "image name/registry search"                           example(ubi8)
// @Param limit           query int      false "limit per page"                                       example(10) minimum(0) maximum(100)
// @Param offset          query int      false "page offset"                                          example(10) minimum(0)
// @Param data_format     query string   false "data section format"                                  enums(json,csv)
// @Param report          query bool     false "overrides limit and offset to return everything"
// @router /clusters/{cluster_id}/exposed_images [get]
// @success 200 {object} GetClusterImagesResponse
// @failure 400 {object} base.Error
// @failure 404 {object} base.Error "cluster does not exist"
// @failure 500 {object} base.Error
func (c *Controller) GetClusterImages(ctx *gin.Context) {
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
	query := c.BuildClusterImagesQuery(accountID, clusterID)

	dataRes := []GetClusterImagesSelect{}
	usedFilters, totalItems, inputErr, dbErr := base.ListQuery(query, getClusterImagesAllowedFilters, filters, getClusterImagesFilterArgs, &dataRes)
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

func (c *Controller) BuildClusterImagesQuery(accountID int64, clusterID uuid.UUID) *gorm.DB {
	return c.Conn.Table("repository").
		Select(`repository.repository, repository.registry, COALESCE(repository_image.tags, '[]') AS tags`).
		Joins("JOIN repository_image ON repository.id = repository_image.repository_id").
		Joins("JOIN cluster_image ON repository_image.image_id = cluster_image.image_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id"). // Do not show images without vulnerabilities
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Where("cluster.account_id = ? AND cluster.uuid = ?", accountID, clusterID)
}
