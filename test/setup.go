package test

import "github.com/gin-gonic/gin"

type Endpoint struct {
	HTTPMethod string
	Path       string
	Handler    gin.HandlerFunc
}

func BuildTestRouter(endpoints []Endpoint, middlewares ...gin.HandlerFunc) *gin.Engine {
	engine := gin.New()

	engine.Use(middlewares...)
	for _, endpoint := range endpoints {
		engine.Handle(endpoint.HTTPMethod, endpoint.Path, endpoint.Handler)
	}
	return engine
}
