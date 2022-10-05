package clusters

import (
	"app/base/ams"
	"app/base/utils"
	"app/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func callGetClusterDetails(t *testing.T, accountID int64, clusterUUID string, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	params := gin.Params{{Key: "cluster_id", Value: clusterUUID}}
	ctx, w := test.MockGinRequest(header, "GET", nil, params, url.Values{})
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetClusterDetails(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetClusterDetailsWrongCluster(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	accID := allAccounts[0].ID

	// Wrong cluster ID causes 400.
	callGetClusterDetails(t, accID, "", http.StatusBadRequest)
}

func TestGetClusterDetails(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		accountClusters := test.GetAccountClusters(t, account.ID)
		for _, cluster := range accountClusters {
			var resp GetClusterDetailsResponse
			w := callGetClusterDetails(t, account.ID, cluster.UUID.String(), http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, cluster.UUID.String(), resp.Data.UUID)
			assert.Equal(t, cluster.LastSeen.UTC(), resp.Data.LastSeen.UTC())
			assert.Equal(t, cluster.UUID.String(), resp.Data.DisplayName)
		}
	}
}

func TestGetClusterDetailsAMS(t *testing.T) {
	utils.Cfg.AmsEnabled = true
	defer func() { utils.Cfg.AmsEnabled = false }()

	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		accountClusters := test.GetAccountClusters(t, account.ID)
		for _, c := range accountClusters {
			testController.AMSClient = &test.AMSClientMock{
				ClusterResponse: ams.ClusterInfo{
					ID:          c.UUID.String(),
					DisplayName: c.UUID.String(),
					Status:      c.Status,
					Type:        c.Type,
					Version:     c.Version,
					Provider:    c.Provider,
				},
			}

			var resp GetClusterDetailsResponse
			w := callGetClusterDetails(t, account.ID, c.UUID.String(), http.StatusOK)

			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, c.UUID.String(), resp.Data.UUID)
			assert.Equal(t, c.LastSeen.UTC(), resp.Data.LastSeen.UTC())
			assert.Equal(t, c.UUID.String(), resp.Data.DisplayName)
		}
	}
}
