package api

import (
	"app/test"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestResp(t *testing.T) {
	strPtr := "test-string-ptr"
	expectedResp := test.Response{
		StringField:    "test-string",
		StringPtrField: &strPtr,
		IntField:       7357,
		Int64Field:     7357,
	}

	expectedCode := 200
	expectedBody, err := json.Marshal(expectedResp)
	assert.Nil(t, err)
	client := NewClientMock("OK", expectedCode, expectedBody)

	var resp test.Response
	c, err := client.Request("GET", "test-url", nil, &resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, c)
	assert.True(t, reflect.DeepEqual(expectedResp, resp))
}

func TestRequestInvalidRespPtr(t *testing.T) {
	client := NewClientMock("Bad request", http.StatusBadRequest, nil)
	_, err := client.Request("GET", "test-url", nil, nil)
	assert.NotNil(t, err)
}

func TestRequestSent(t *testing.T) {
	expectedURL := "test-url"
	expectedMethod := "GET"
	expectedHeaders := map[string][]string{"Content-Type": {"application/json"}}

	expectedBody := test.Request{StringField: "test-string"}
	bodyBytes, err := json.Marshal(test.Request{StringField: "test-string"})
	bodyBytes = append(bodyBytes, 0xA)
	assert.Nil(t, err)

	client := NewClientRequestCheckerMock(t, expectedMethod, expectedURL, expectedHeaders, bodyBytes)
	_, err = client.Request(expectedMethod, expectedURL, &expectedBody, &test.Response{})
	assert.Nil(t, err)
}

func TestRetryRequest(t *testing.T) {
	strPtr := "test-string-ptr"
	expectedResp := test.Response{
		StringField:    "test-string",
		StringPtrField: &strPtr,
		IntField:       7357,
		Int64Field:     7357,
	}

	expectedCode := 200
	expectedBody, err := json.Marshal(expectedResp)
	assert.Nil(t, err)
	client := NewClientMock("OK", expectedCode, expectedBody)

	var resp test.Response
	c, err := client.RetryRequest("GET", "test-url", nil, &resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedCode, c)
	assert.True(t, reflect.DeepEqual(expectedResp, resp))
}

func TestRetryRequestFail(t *testing.T) {
	retries = 5

	client := NewClientMock("Bad request", http.StatusBadRequest, nil)
	c, err := client.RetryRequest("GET", "test-url", nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, c, http.StatusBadRequest)
}

func NewClientMock(status string, code int, resp []byte) Client {
	return Client{
		HTTPClient: &test.HTTPClientMock{
			Status:     status,
			StatusCode: code,
			RespBytes:  resp,
		}}
}

func NewClientRequestCheckerMock(t *testing.T, method, url string, headers map[string][]string, body []byte) Client {
	return Client{
		HTTPClient: &test.HTTPRequestChecker{
			ExpectedMethod:  method,
			ExpectedURL:     url,
			ExpectedHeaders: headers,
			ExpectedBody:    body,
			T:               t,
		},
	}
}
