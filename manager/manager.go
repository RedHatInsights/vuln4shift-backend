package manager

import (
	"app/manager/controllers/meta"
	"app/manager/middlewares"

	"github.com/gin-gonic/gin"
)

// createMetaGroup adds meta endpoints to the router.
func createMetaGroup(router *gin.Engine) *gin.RouterGroup {
	metaGroup := router.Group("/")

	metaGroup.GET("healthz", meta.GetApistatus)
	metaGroup.GET("apistatus", meta.GetApistatus)
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

// Start starts manager router.
func Start() {
	router := BuildRouter()
	err := router.Run(":8000")

	if err != nil {
		panic(err)
	}
}
