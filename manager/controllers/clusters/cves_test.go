package clusters

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

func callGetClusterCves(t *testing.T, accountID int64, clusterUUID string, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	params := gin.Params{{Key: "cluster_id", Value: clusterUUID}}
	ctx, w := test.MockGinRequest(header, "GET", nil, params, url.Values{})
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetClusterCves(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetClusterCvesWrongCluster(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	accID := allAccounts[0].ID

	// Wrong cluster ID causes 400.
	callGetClusterCves(t, accID, "", http.StatusBadRequest)
}

func TestGetClusterCves(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		accountClusters := test.GetAccountClusters(t, account.ID)
		for _, cluster := range accountClusters {
			expectedCves := test.GetClusterCves(t, cluster.ID)
			var resp GetClusterCvesResponse
			w := callGetClusterCves(t, account.ID, cluster.UUID.String(), http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, len(expectedCves), len(resp.Data))
			for i, ec := range expectedCves {
				ac := resp.Data[i]
				assert.Equal(t, ec.Name, *ac.Name)
				assert.Equal(t, ec.Cvss2Score, *ac.Cvss2Score)
				assert.Equal(t, ec.Cvss3Score, *ac.Cvss3Score)
				assert.Equal(t, ec.Description, *ac.Description)
				assert.Equal(t, ec.Severity, *ac.Severity)
				assert.Equal(t, ec.PublicDate, ac.PublicDate)
			}
			totalItems := test.GetMetaTotalItems(resp.Meta)
			assert.Equal(t, float64(len(resp.Data)), totalItems)
		}
	}
}
