package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type HTTPClientMock struct {
	Status     string
	StatusCode int
	RespBytes  []byte
}

func NewAPIMock(status string, code int, response []byte) *HTTPClientMock {
	return &HTTPClientMock{
		Status:     status,
		StatusCode: code,
		RespBytes:  response,
	}
}

func (m *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader(m.RespBytes))
	return &http.Response{
		Status:     m.Status,
		StatusCode: m.StatusCode,
		Body:       r,
	}, nil
}
