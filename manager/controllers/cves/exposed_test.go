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

func callGetExposedClusters(t *testing.T, accountID int64, cveName string, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	params := gin.Params{{Key: "cve_name", Value: cveName}}

	ctx, w := test.MockGinRequest(header, "GET", nil, params, url.Values{})
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetExposedClusters(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetExposedClusters(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		cves := test.GetAccountsCves(t, account.ID)
		for _, cve := range cves {
			var resp GetExposedClustersResponse
			w := callGetExposedClusters(t, account.ID, cve.Name, http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

			expectedClusters := test.GetExposedClusters(t, account.ID, cve.ID)
			assert.Equal(t, len(expectedClusters), len(resp.Data))
			for i, ec := range expectedClusters {
				// Actual cluster.
				ac := resp.Data[i]
				assert.Equal(t, ec.UUID.String(), ac.DisplayName)
				assert.Equal(t, ec.UUID.String(), ac.UUID)
				assert.Equal(t, ec.LastSeen.UTC(), test.GetUTC(ac.LastSeen))
				assert.Equal(t, ec.Status, ac.Status)
				assert.Equal(t, ec.Version, ac.Version)
				assert.Equal(t, ec.Provider, ac.Provider)
			}
		}
	}
}
