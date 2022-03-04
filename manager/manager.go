package manager

import (
	"app/manager/controllers/meta"
	"app/manager/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var openAPILocation = "/api/vuln4shift/v1/openapi.json"

// createMetaGroup adds meta endpoints to the router.
func createMetaGroup(router *gin.Engine) *gin.RouterGroup {
	metaGroup := router.Group("/")

	metaGroup.GET("healthz", meta.GetApistatus)
	metaGroup.GET("apistatus", meta.GetApistatus)

	openAPIURL := ginSwagger.URL(openAPILocation)
	metaGroup.GET("openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, openAPIURL))
	metaGroup.StaticFile(openAPILocation, "./manager/docs/swagger.json")
	return metaGroup
}

// setMiddlewares sets middlewares for router.
func setMiddlewares(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(middlewares.Logger())
}

// BuildRouter creates manager router with endpoints and middlewares.
func BuildRouter() *gin.Engine {
	router := gin.New()

	setMiddlewares(router)
	createMetaGroup(router)
	return router
}

// @title Vulnerability for Openshift API documentation
// @version 0.1.0
// @description Documentation to the REST API for application
// @description Vulnerability for Openshift based on console.redhat.com.
//
// @securityDefinitions.apikey RhIdentity
// @in header
// @name x-rh-identity
//
// @query.collection.format multi
// @basepath /
// @schemes http
func Start() {
	router := BuildRouter()
	err := router.Run(":8000")

	if err != nil {
		panic(err)
	}
}
