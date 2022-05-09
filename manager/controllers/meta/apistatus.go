package meta

import "github.com/gin-gonic/gin"

func (e *Controller) GetApistatus(ctx *gin.Context) {
	ctx.Status(200)
}
