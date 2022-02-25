package meta

import "github.com/gin-gonic/gin"

// GetApistatus represents health/status endpoint controller.
func GetApistatus(ctx *gin.Context) {
	ctx.Status(200)
}
