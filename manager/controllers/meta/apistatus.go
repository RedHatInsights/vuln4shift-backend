package meta

import "github.com/gin-gonic/gin"

// GetApistatus represents health/status endpoint controller.
//
// @id GetApistatus
// @summary API status of the application
// @description Endpoint checks for status of the application
// @accept */*
// @router /apistatus [get]
// @success 200
// @failure 503
func (e *Controller) GetApistatus(ctx *gin.Context) {
	ctx.Status(200)
}
