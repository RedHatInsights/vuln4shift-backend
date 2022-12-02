package cves

import (
	"app/base/ams"
	"app/base/logging"
	"app/base/models"
	"app/base/utils"
	"app/manager/middlewares"
	"app/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testFilterer   = middlewares.Filterer()
	testController Controller
)

func callGetCves(t *testing.T, accountID int64, affectedClusters bool, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	urlValues := url.Values{}
	if affectedClusters {
		urlValues["affected_clusters"] = []string{"true", "false"}
	}

	ctx, w := test.MockGinRequest(header, "GET", nil, gin.Params{}, urlValues)
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetCves(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetCvesNonAffecting(t *testing.T) {
	account := test.GetAccounts(t)[0]

	var resp GetCvesResponse
	w := callGetCves(t, account.ID, false, http.StatusOK)
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

	allCves := test.GetAllCves(t)
	assert.Equal(t, len(allCves), len(resp.Data))
}

func TestGetCvesAffecting(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		var resp GetCvesResponse
		w := callGetCves(t, account.ID, true, http.StatusOK)
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

		expectedCves := test.GetAccountCves(t, account.ID)
		assert.Equal(t, len(resp.Data), len(expectedCves))
		for i, ec := range expectedCves {
			// Actual CVE
			ac := resp.Data[i]
			assert.Equal(t, ec.Name, test.GetStringPtrValue(ac.Name))
			assert.Equal(t, ec.Severity, *ac.Severity)
			assert.Equal(t, ec.Cvss2Score, test.GetFloat32PtrValue(ac.Cvss2Score))
			assert.Equal(t, ec.Description, test.GetStringPtrValue(ac.Description))
			// assert.Equal(t, test.GetImagesExposed(t, account.ID, ec.ID), *ac.ImagesExposed)
			assert.Equal(t, test.GetClustersExposed(t, account.ID, ec.ID), *ac.ClustersExposed)
		}
		totalItems := test.GetMetaTotalItems(resp.Meta)
		assert.Equal(t, float64(len(resp.Data)), totalItems)
	}
}

func TestGetCvesAffectingAMS(t *testing.T) {
	utils.Cfg.AmsEnabled = true
	defer func() { utils.Cfg.AmsEnabled = false }()

	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		allClusters := test.GetAccountClusters(t, account.ID)
		// Some subset of the account clusters
		expectedClusters := allClusters[:len(allClusters)/2]

		amsClusters := make(map[string]ams.ClusterInfo)
		for _, c := range expectedClusters {
			amsClusters[c.UUID.String()] = ams.ClusterInfo{
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

		var resp GetCvesResponse
		w := callGetCves(t, account.ID, true, http.StatusOK)
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

		clusterIDs := make([]int64, 0, len(expectedClusters))
		for _, c := range expectedClusters {
			clusterIDs = append(clusterIDs, c.ID)
		}

		expectedCves := test.GetAccountCvesForClusters(t, account.ID, clusterIDs)
		assert.Equal(t, len(expectedCves), len(resp.Data))
		for i, ec := range expectedCves {
			// Actual CVE
			ac := resp.Data[i]
			assert.Equal(t, ec.Name, test.GetStringPtrValue(ac.Name))
			assert.Equal(t, ec.Severity, *ac.Severity)
			assert.Equal(t, ec.Cvss2Score, test.GetFloat32PtrValue(ac.Cvss2Score))
			assert.Equal(t, ec.Description, test.GetStringPtrValue(ac.Description))
			// assert.Equal(t, test.GetImagesExposed(t, account.ID, ec.ID), *ac.ImagesExposed)
			assert.Equal(t, test.GetClustersExposedLimitClusters(t, account.ID, ec.ID, clusterIDs), *ac.ClustersExposed)
		}
		totalItems := test.GetMetaTotalItems(resp.Meta)
		assert.Equal(t, float64(len(resp.Data)), totalItems)
	}
}

func TestMain(m *testing.M) {
	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		panic(err)
	}

	test.DB = db
	err = test.ResetDB()
	if err != nil {
		panic(err)
	}

	logger, err := logging.CreateLogger("DEBUG")
	if err != nil {
		panic(err)
	}

	testController.Conn = test.DB
	testController.Logger = logger

	os.Exit(m.Run())
}
