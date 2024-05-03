package cves

import (
	"app/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func callGetCveImages(t *testing.T, accountID int64, cveName string,
	expectedStatus int, filters map[string][]string) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	params := gin.Params{{Key: "cve_name", Value: cveName}}

	urlValues := url.Values{}
	if filters != nil {
		for filter, value := range filters {
			urlValues[filter] = value
		}
	}

	ctx, w := test.MockGinRequest(header, "GET", nil, params, urlValues)
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetCveImages(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}
func TestGetCveImagesWrongCve(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	accID := allAccounts[0].ID

	// Empty cve name causes 400.
	callGetCveImages(t, accID, "", http.StatusNotFound, nil)
}

func TestGetCveImages(t *testing.T) {
	accID1 := 13
	cveName1 := "CVE-2022-0001"
	clustersExposed := int32(1)
	repository := "rhel7.1"
	registry := "registry.access.redhat.com"
	version := "Unknown"

	var resp1 GetCveImagesResponse
	w := callGetCveImages(t, int64(accID1), cveName1, http.StatusOK, nil)
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp1))
	assert.Equal(t, 1, len(resp1.Data))
	assert.Equal(t, registry, *resp1.Data[0].Registry)
	assert.Equal(t, repository, *resp1.Data[0].Repository)
	assert.Equal(t, version, *resp1.Data[0].Version)
	assert.Equal(t, clustersExposed, *resp1.Data[0].ClustersExposed)
}
