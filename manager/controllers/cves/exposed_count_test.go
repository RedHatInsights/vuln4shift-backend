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

func callGetExposedClustersCount(t *testing.T, accountID int64, cveName string, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	params := gin.Params{{Key: "cve_name", Value: cveName}}

	ctx, w := test.MockGinRequest(header, "GET", nil, params, url.Values{})
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetExposedClustersCount(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetExposedClustersCount(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		cves := test.GetAccountCves(t, account.ID)
		for _, cve := range cves {
			var resp GetExposedClustersCountResponse
			w := callGetExposedClustersCount(t, account.ID, cve.Name, http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

			testController.Logger.Infof("account id=%v, cve_id=%v", account.ID, cve.ID)
			expectedClusters := test.GetExposedClusters(t, account.ID, cve.ID)
			assert.Equal(t, int64(len(expectedClusters)), resp.Count)
		}
	}
}
