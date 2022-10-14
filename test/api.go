package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Request struct {
	StringField    string
	StringPtrField *string
	IntField       int
	Int64Field     int64
}

type Response struct {
	StringField    string
	StringPtrField *string
	IntField       int
	Int64Field     int64
}

type HTTPClientMock struct {
	Status     string
	StatusCode int
	RespBytes  []byte
}

func (m *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader(m.RespBytes))
	return &http.Response{
		Status:     m.Status,
		StatusCode: m.StatusCode,
		Body:       r,
	}, nil
}

type HTTPRequestChecker struct {
	ExpectedMethod  string
	ExpectedURL     string
	ExpectedHeaders map[string][]string
	ExpectedBody    []byte
	T               *testing.T
}

func (c *HTTPRequestChecker) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(c.T, c.ExpectedMethod, req.Method)
	assert.Equal(c.T, c.ExpectedURL, req.URL.Path)

	for key, expectedVal := range c.ExpectedHeaders {
		actualVal, found := req.Header[key]
		assert.True(c.T, found)

		sort.Strings(expectedVal)
		sort.Strings(actualVal)

		assert.Equal(c.T, expectedVal, actualVal)
	}

	reqBytes, err := ioutil.ReadAll(req.Body)
	assert.Nil(c.T, err)
	assert.Equal(c.T, c.ExpectedBody, reqBytes)

	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(c.ExpectedBody)),
	}, nil
}
