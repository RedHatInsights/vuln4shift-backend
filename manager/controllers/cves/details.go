package cves

import (
	"app/base/models"
	"app/manager/base"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCveDetailsSelect
// @Description CVE details data
// @Description presents in response
type GetCveDetailsSelect struct {
	Cvss2Score   *float32        `json:"cvss2_score"`
	Cvss2Metrics *string         `json:"cvss2_metrics"`
	Cvss3Score   *float32        `json:"cvss3_score"`
	Cvss3Metrics *string         `json:"cvss3_metrics"`
	Description  string          `json:"description"`
	Severity     models.Severity `json:"severity"`
	PublicDate   *time.Time      `json:"publish_date"`
	Name         string          `json:"synopsis"`
	RedhatURL    *string         `json:"redhat_url"`
}

type GetCveDetailsResponse struct {
	Data GetCveDetailsSelect `json:"data"`
	Meta interface{}         `json:"meta"`
}

// GetCveDetails represents CVE detail endpoint controller.
//
// @id GetCveDetails
// @summary CVE details
// @security RhIdentity || BasicAuth
// @Tags cves
// @description Endpoint return details for given CVE
// @accept */*
// @produce json
// @Param cve_name path  string true  "CVE name"
// @router /cves/{cve_name} [get]
// @success 200 {object} GetCveDetailsResponse
// @failure 404 {object} base.Error "{cve_name} not found"
// @failure 500 {object} base.Error
func (c *Controller) GetCveDetails(ctx *gin.Context) {
	cveName := ctx.Param("cve_name")
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

	ctx.JSON(
		http.StatusOK,
		GetCveDetailsResponse{cveDetails, base.BuildMeta(make(map[string]base.Filter), nil, nil, nil, nil)},
	)
}

func (c *Controller) BuildCveDetailsQuery(cveName string) *gorm.DB {
	return c.Conn.Table("cve").
		Select(`name, description, public_date, severity, cvss2_score, cvss2_metrics, 
                       cvss3_score, cvss3_metrics, redhat_url`).
		Where("cve.name = ?", cveName)
}
