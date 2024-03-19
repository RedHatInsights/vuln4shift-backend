package clusters

import (
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

func callGetClusterImages(t *testing.T, accountID int64, clusterUUID string,
	expectedStatus int, filters map[string][]string) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	params := gin.Params{{Key: "cluster_id", Value: clusterUUID}}

	urlValues := url.Values{}
	if filters != nil {
		for filter, value := range filters {
			urlValues[filter] = value
		}
	}

	ctx, w := test.MockGinRequest(header, "GET", nil, params, urlValues)
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetClusterImages(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetClusterImagesWrongCluster(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	accID := allAccounts[0].ID

	// Wrong cluster ID causes 400.
	callGetClusterImages(t, accID, "", http.StatusBadRequest, nil)
}

func TestGetClusterImages(t *testing.T) {
	allAccounts := test.GetAccounts(t)

	for _, account := range allAccounts {
		accountClusters := test.GetAccountClusters(t, account.ID)

		for _, cluster := range accountClusters {
			expectedRepoImages := test.GetClusterRepoImages(t, cluster.ID)
			var resp GetClusterImagesResponse

			w := callGetClusterImages(t, account.ID, cluster.UUID.String(), http.StatusOK, nil)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

			assert.Equal(t, len(expectedRepoImages), len(resp.Data))

			for i, eri := range expectedRepoImages {
				er := test.GetRepoByID(t, eri.RepositoryID)
				var expectedImageVersion utils.ImageVersion
				if eri.Tags == nil {
					expectedImageVersion = utils.Unknown
				} else {
					tags := string(eri.Tags.Bytes)
					_ = expectedImageVersion.Scan(tags)
				}
				ac := resp.Data[i]
				assert.Equal(t, er.Repository, *ac.Repository)
				assert.Equal(t, er.Registry, *ac.Registry)
				assert.Equal(t, expectedImageVersion, *ac.Version)
			}

			totalItems := test.GetMetaTotalItems(resp.Meta)
			assert.Equal(t, float64(len(resp.Data)), totalItems)
		}
	}
}

func TestGetClusterImagesVulnerableOnly(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		accountClusters := test.GetAccountClusters(t, account.ID)

		for _, cluster := range accountClusters {
			var resp GetClusterImagesResponse

			w := callGetClusterImages(t, account.ID, cluster.UUID.String(), http.StatusOK, nil)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

			for _, ac := range resp.Data {
				assert.NotEqual(t, *ac.Repository, "rhel6.10.novulns")
			}
		}
	}
}
