package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetMetricsServer(t *testing.T, method, job string) *httptest.Server {
	testURI := fmt.Sprintf("/metrics/job/%s", job)

	checkPush := func(_ http.ResponseWriter, r *http.Request) {
		assert.Equal(t, method, r.Method)
		assert.Equal(t, testURI, r.RequestURI)
	}

	return httptest.NewServer(http.HandlerFunc(checkPush))
}
