package cves

import (
	"app/base/logging"
	"app/base/models"
	"app/base/utils"
	"app/manager/base"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

// GetCveDetails represents CVE detail endpoint controller.
//
// @id GetExposedClusters
// @summary CVE details
// @Tags cves
// @description Endpoint return details for given CVE
// @accept */*
// @produce json
// @Param cve_name path  string true  "CVE name"
// @router /cves/{cve_name} [get]
// @success 200 {array}  base.Response{data=GetCveDetailsSelect}
// @failure 404 {object} base.Error{error=base.ErrorDetail{detail=string,status=int}} "{cve_name} not found"
// @failure 500 {object} base.Error{error=base.ErrorDetail{detail=string,status=int}} "Internal server error
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
		logger.Errorf("Database error: %s", result.Error)
		return
	}

	ctx.JSON(
		http.StatusOK,
		base.BuildResponse(cveDetails, base.BuildMeta(make(map[string]base.Filter), make([]string, 0))),
	)
}

func (c *Controller) BuildCveDetailsQuery(cveName string) *gorm.DB {
	return c.Conn.Table("cve").
		Select(`name, description, public_date, severity, cvss2_score, cvss2_metrics, 
                       cvss3_score, cvss3_metrics, redhat_url`).
		Where("cve.name = ?", cveName)
}
