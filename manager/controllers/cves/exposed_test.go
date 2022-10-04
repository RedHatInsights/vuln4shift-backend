package cves

import (
	"app/base/utils"
	"app/manager/amsclient"
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
		cves := test.GetAccountCves(t, account.ID)
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

func TestGetExposedClustersAMSSubset(t *testing.T) {
	utils.Cfg.AmsEnabled = true
	defer func() { utils.Cfg.AmsEnabled = false }()

	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		cves := test.GetAccountCves(t, account.ID)
		for _, cve := range cves {
			expectedClusters := test.GetExposedClusters(t, account.ID, cve.ID)
			if len(expectedClusters) == 0 {
				continue
			}

			c := expectedClusters[0]
			testController.AMSClient = &test.AMSClientMock{
				ClustersResponse: map[string]amsclient.ClusterInfo{
					c.UUID.String(): {
						ID:          c.UUID.String(),
						DisplayName: c.UUID.String(),
						Status:      c.Status,
						Type:        c.Type,
						Version:     c.Version,
						Provider:    c.Provider,
					},
				},
			}

			var resp GetExposedClustersResponse
			w := callGetExposedClusters(t, account.ID, cve.Name, http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, 1, len(resp.Data))

			ec := expectedClusters[0]
			ac := resp.Data[0]
			assert.Equal(t, ec.UUID.String(), ac.DisplayName)
			assert.Equal(t, ec.UUID.String(), ac.UUID)
			assert.Equal(t, ec.LastSeen.UTC(), test.GetUTC(ac.LastSeen))
			assert.Equal(t, ec.Status, ac.Status)
			assert.Equal(t, ec.Version, ac.Version)
			assert.Equal(t, ec.Provider, ac.Provider)
		}
	}
}

func TestGetExposedClustersAMS(t *testing.T) {
	utils.Cfg.AmsEnabled = true
	defer func() { utils.Cfg.AmsEnabled = false }()

	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		cves := test.GetAccountCves(t, account.ID)
		for _, cve := range cves {
			expectedClusters := test.GetExposedClusters(t, account.ID, cve.ID)

			expectedProviders := make(map[string]bool)
			expectedStatuses := make(map[string]bool)
			expectedVersions := make(map[string]bool)

			amsClusters := make(map[string]amsclient.ClusterInfo)
			for _, c := range expectedClusters {
				expectedProviders[c.Provider] = true
				expectedStatuses[c.Status] = true
				expectedVersions[c.Version] = true
				amsClusters[c.UUID.String()] = amsclient.ClusterInfo{
					ID:          c.UUID.String(),
					DisplayName: c.UUID.String(),
					Status:      c.Status,
					Type:        c.Type,
					Version:     c.Version,
					Provider:    c.Provider,
				}
			}
			testController.AMSClient = &test.AMSClientMock{
				ClustersResponse: amsClusters,
			}

			var resp GetExposedClustersResponse
			w := callGetExposedClusters(t, account.ID, cve.Name, http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
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

			if len(resp.Data) > 0 {
				test.CheckClustersMeta(t, resp.Meta, expectedProviders, expectedStatuses, expectedVersions)
			}
		}
	}
}
