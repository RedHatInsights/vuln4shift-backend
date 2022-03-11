package meta

import (
	"app/manager/middlewares"
	"app/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func buildApistatusRouter() *gin.Engine {
	metaEndpoint := Controller{
		Conn: nil,
	}

	endpoints := []test.Endpoint{
		{
			HTTPMethod: "GET",
			Path:       "/apistatus",
			Handler:    metaEndpoint.GetApistatus,
		},
	}

	return test.BuildTestRouter(endpoints, middlewares.Logger())
}

func TestApistatus(t *testing.T) {
	router := buildApistatusRouter()

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/apistatus", nil)

	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code, "Status API should respond 200 code")
}
