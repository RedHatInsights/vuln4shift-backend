package manager

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"app/base/models"
	"app/base/utils"
	"app/manager/controllers/cves"
	"app/manager/controllers/meta"
	"app/manager/middlewares"
)

var (
	apiPrefix       = "/api/vuln4shift"
	openAPILocation = fmt.Sprintf("%s/v1/openapi.json", apiPrefix)
)

// createMetaGroup adds meta endpoints to the router.
func createMetaGroup(router *gin.Engine, db *gorm.DB) *gin.RouterGroup {
	metaGroup := router.Group("/")

	metaController := meta.Controller{
		Conn: db,
	}

	metaGroup.GET("healthz", metaController.GetApistatus)
	metaGroup.GET("apistatus", metaController.GetApistatus)

	openAPIURL := ginSwagger.URL(openAPILocation)
	metaGroup.GET(fmt.Sprintf("%s/v1/openapi/*any", apiPrefix), ginSwagger.WrapHandler(swaggerFiles.Handler, openAPIURL))
	metaGroup.StaticFile(openAPILocation, "./manager/docs/swagger.json")
	return metaGroup
}

func createCveGroup(router *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	cveGroup := router.Group("/v1/cves")

	cveController := cves.Controller{
		Conn: db,
	}

	// Cves endpoints must be authenticated
	cveGroup.Use(middlewares.Authenticate(db))

	cveGroup.GET("", cveController.GetCves)
	cveGroup.GET("/:cve_name/exposed_clusters", cveController.GetExposedClusters)
	cveGroup.GET("/:cve_name", cveController.GetCveDetails)
	return cveGroup
}

// setMiddlewares sets middlewares for router.
func setMiddlewares(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(middlewares.Logger())
	router.Use(middlewares.Filterer())
}

// BuildRouter creates manager router with endpoints and middlewares.
func BuildRouter() *gin.Engine {
	router := gin.New()

	dsn := utils.GetDbURL(false)
	db, err := models.GetGormConnection(dsn)

	if err != nil {
		log.Fatalf(err.Error())
	}

	setMiddlewares(router)
	createMetaGroup(router, db)

	api := router.Group(apiPrefix)
	createCveGroup(api, db)

	return router
}

// Start
//
// @title Vulnerability for Openshift API documentation
// @version 0.1.0
// @description Documentation to the REST API for application
// @description Vulnerability for Openshift based on console.redhat.com.
//
// @securityDefinitions.apikey RhIdentity
// @in header
// @name x-rh-identity
//
// @securityDefinitions.basic BasicAuth
//
// @query.collection.format multi
// @basePath /api/vuln4shift/v1
// @schemes http https
func Start() {
	router := BuildRouter()
	err := router.Run(fmt.Sprintf(":%d", utils.Cfg.PublicPort))

	if err != nil {
		panic(err)
	}
}
