package middlewares

import (
	"app/base/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger represents Logger middleware factory.
// Function sets up a logger which is later used for every
// request on the API to be logged.
func Logger() gin.HandlerFunc {
	logger, err := utils.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		panic("Invalid LOGGING_LEVEL environment variable set")
	}

	return func(ctx *gin.Context) {
		startTimestamp := time.Now()
		ctx.Next()
		duration := time.Since(startTimestamp)

		entry := logger.WithFields(logrus.Fields{
			"timestamp": startTimestamp.UTC().Format(time.RFC3339),
			"method":    ctx.Request.Method,
			"status":    ctx.Writer.Status(),
			"duration":  duration.String(),
			"org_id":    ctx.GetString("org_id"),
		})
		if ctx.Writer.Status() < http.StatusInternalServerError {
			entry.Info(ctx.Request.RequestURI)
		} else {
			entry.Error(ctx.Request.RequestURI)
		}
	}
}
