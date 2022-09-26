package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

func MockGinRequest(header http.Header, method string, body []byte, params gin.Params, urlValues url.Values) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = &http.Request{
		Body:   io.NopCloser(bytes.NewBuffer(body)),
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	ctx.Request.Method = method
	ctx.Request.Header = header
	ctx.Params = params
	ctx.Request.URL.RawQuery = urlValues.Encode()
	return ctx, r
}
