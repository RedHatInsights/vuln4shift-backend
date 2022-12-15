package amsclient

import (
	"app/base/ams"
	"app/base/models"
	"app/base/utils"
	"app/test"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	sdk "github.com/openshift-online/ocm-sdk-go"

	"github.com/stretchr/testify/assert"
)

func prepareTestClusterDetails(t *testing.T, accountID int64) ([]models.Cluster, map[string]bool, map[string]bool, map[string]bool, map[string]ams.ClusterInfo) {
	expectedClusters := test.GetAccountClusters(t, accountID)

	expectedProviders := make(map[string]bool)
	expectedStatuses := make(map[string]bool)
	expectedVersions := make(map[string]bool)

	amsClusters := make(map[string]ams.ClusterInfo)
	for _, c := range expectedClusters {
		expectedProviders[c.Provider] = true
		expectedStatuses[c.Status] = true
		expectedVersions[c.Version] = true

		amsClusters[c.UUID.String()] = ams.ClusterInfo{
			ID:          c.UUID.String(),
			DisplayName: c.UUID.String(),
			Status:      c.Status,
			Type:        c.Type,
			Version:     c.Version,
			Provider:    c.Provider,
		}
	}

	return expectedClusters, expectedProviders, expectedStatuses, expectedVersions, amsClusters
}

func TestDBFetchClusterDetails(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		expectedClusters, expectedProviders, expectedStatuses, expectedVersions, amsClusters := prepareTestClusterDetails(t, account.ID)

		amsMock := &test.AMSClientMock{
			ClustersResponse: amsClusters,
		}

		actualUUIDs, statuses, versions, providers, err := DBFetchClusterDetails(test.DB, amsMock, account.ID, account.OrgID, true, nil)
		assert.Nil(t, err)

		assert.Equal(t, len(expectedClusters), len(actualUUIDs))
		for i, uuid := range actualUUIDs {
			assert.Equal(t, uuid, expectedClusters[i].UUID.String())
		}

		if len(actualUUIDs) > 0 {
			test.CheckClusterDetails(t, expectedProviders, expectedStatuses, expectedVersions, providers, statuses, versions)
		}
	}
}

func TestDBFetchClusterDetailsCve(t *testing.T) {
	account := test.GetAccounts(t)[1]

	expectedClusters, _, _, _, amsClusters := prepareTestClusterDetails(t, account.ID)

	cluster := expectedClusters[0]

	amsMock := &test.AMSClientMock{
		ClustersResponse: amsClusters,
	}

	// Test against random CVE
	cves := test.GetClusterCves(t, cluster.ID)
	if len(cves) == 0 {
		t.Skip()
	}
	cve := cves[0]

	actualUUIDs, _, _, _, err := DBFetchClusterDetails(test.DB, amsMock, account.ID, account.OrgID, true, &cve.Name)
	assert.Nil(t, err)

	var clusterFound bool
	for _, uuid := range actualUUIDs {
		if uuid == cluster.UUID.String() {
			clusterFound = true
		}
	}

	assert.True(t, clusterFound)
}

func TestDBFetchClusterDetailsCveDoesNotExist(t *testing.T) {
	account := test.GetAccounts(t)[0]

	_, _, _, _, amsClusters := prepareTestClusterDetails(t, account.ID)
	amsMock := &test.AMSClientMock{
		ClustersResponse: amsClusters,
	}

	cveName := "does-not-exist"

	actualUUIDs, _, _, _, err := DBFetchClusterDetails(test.DB, amsMock, account.ID, account.OrgID, true, &cveName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(actualUUIDs))
}

func TestDBFetchClusterDetailsSyncFalse(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		expectedClusters := test.GetAccountClusters(t, account.ID)

		expectedProviders := make(map[string]bool)
		expectedStatuses := make(map[string]bool)
		expectedVersions := make(map[string]bool)

		for _, c := range expectedClusters {
			expectedProviders[c.Provider] = true
			expectedStatuses[c.Status] = true
			expectedVersions[c.Version] = true
		}

		uuids, statuses, versions, providers, err := DBFetchClusterDetails(test.DB, nil, account.ID, account.OrgID, false, nil)
		assert.Nil(t, err)

		for i, uuid := range uuids {
			assert.Equal(t, uuid, expectedClusters[i].UUID.String())
		}

		if len(uuids) > 0 {
			test.CheckClusterDetails(t, expectedProviders, expectedStatuses, expectedVersions, providers, statuses, versions)
		}
	}
}

func TestDBFetchClusterDetailsSyncUpdate(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		expectedClusters := test.GetAccountClusters(t, account.ID)

		name := "new-display-name"
		st := "new-status"
		tp := "new-type"
		ver := "new-version"
		prov := "new-provider"

		amsClusters := make(map[string]ams.ClusterInfo)
		for _, c := range expectedClusters {
			amsClusters[c.UUID.String()] = ams.ClusterInfo{
				ID:          c.UUID.String(),
				DisplayName: name,
				Status:      st,
				Type:        tp,
				Version:     ver,
				Provider:    prov,
			}
		}
		amsMock := &test.AMSClientMock{
			ClustersResponse: amsClusters,
		}

		uuids, statuses, versions, providers, err := DBFetchClusterDetails(test.DB, amsMock, account.ID, account.OrgID, true, nil)
		assert.Nil(t, err)

		if len(uuids) == 0 {
			continue
		}

		for i, uuid := range uuids {
			assert.Equal(t, uuid, expectedClusters[i].UUID.String())
		}

		expectedStatuses := map[string]bool{st: true}
		expectedVersions := map[string]bool{ver: true}
		expectedProviders := map[string]bool{prov: true}
		test.CheckClusterDetails(t, expectedProviders, expectedStatuses, expectedVersions, providers, statuses, versions)

		for _, ec := range expectedClusters {
			ac := test.GetCluster(t, ec.ID)
			assert.Equal(t, name, ac.DisplayName)
			assert.Equal(t, st, ac.Status)
			assert.Equal(t, tp, ac.Type)
			assert.Equal(t, ver, ac.Version)
			assert.Equal(t, prov, ac.Provider)
		}
	}
}

func TestGenerateSearchParameter(t *testing.T) {
	acc := test.GetAccounts(t)[0]
	res := generateSearchParameter(acc.OrgID, []string{"ready", "not ready"})
	assert.Equal(t, fmt.Sprintf("organization_id is '%s' and cluster_id != '' and status not in ('ready','not ready')", acc.OrgID), res)
}

func TestNewAMSClient(t *testing.T) {
	utils.Cfg.AmsClientID = "7357"
	utils.Cfg.AmsClientSecret = "7357"

	amsc, err := NewAMSClient()
	assert.Nil(t, err)
	assert.NotNil(t, amsc)
}

type TestAMSTransport struct {
	APIResponses map[string][]byte
	PageLimit    int
}

func (t *TestAMSTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	hdr := http.Header{}
	hdr.Add("content-type", "application/json")

	// Escape infinite loop if applicable.
	if t.PageLimit > 0 {
		page, err := strconv.Atoi(req.URL.Query().Get("page"))
		if err == nil && page > t.PageLimit {
			return &http.Response{StatusCode: 204, Header: hdr}, nil
		}
	}

	return &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Header:     hdr,
		Body:       io.NopCloser(bytes.NewBuffer(t.APIResponses[req.URL.Path])),
	}, nil
}

func newTestAMSClient(t *testing.T, APIResponses map[string][]byte, pageLimit int) *amsClientImpl {
	utils.Cfg.AmsClientID = "test"
	utils.Cfg.AmsClientSecret = "test"
	utils.Cfg.AmsAPIURL = "localhost:7357"

	builder := sdk.NewConnectionBuilder().URL(fmt.Sprintf("http://%s", utils.Cfg.AmsAPIURL))
	builder.TransportWrapper(func(http.RoundTripper) http.RoundTripper {
		return &TestAMSTransport{APIResponses: APIResponses, PageLimit: pageLimit}
	})
	builder = builder.Client(utils.Cfg.AmsClientID, utils.Cfg.AmsClientSecret)

	builder = builder.TokenURL(fmt.Sprintf("http://%s/token", utils.Cfg.AmsAPIURL))
	conn, err := builder.Build()
	assert.Nil(t, err)

	logger, err := utils.CreateLogger(utils.Cfg.LoggingLevel)
	assert.Nil(t, err)

	return &amsClientImpl{
		connection: conn,
		pageSize:   utils.Cfg.AmsAPIPagesize,
		logger:     logger,
	}
}

const (
	tokenAPIResp = `{
	  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	}`

	orgAPIResp = `{
	  "items": [
		{
		  "id": "internal-test-id",
		  "external_id": "external-test-id",
			"name": "test-org"
		}
	  ],
	  "size": 1
	}`

	subAPIResp = `{
	  "status": 200,
	  "items": [
		{
		  "id": "test-sub",
		  "cluster_id": "test-cluster-id",
          "external_cluster_id": "test-external-cluster-id",
          "status": "archived",
          "region_id": "3",
          "display_name": "test-display-name"
		}
	  ],
	  "size": 1
	}`
)

func TestGetClustersForOrganization(t *testing.T) {
	APIResponses := map[string][]byte{
		"/token":                              []byte(tokenAPIResp),
		"/api/accounts_mgmt/v1/organizations": []byte(orgAPIResp),
		"/api/accounts_mgmt/v1/subscriptions": []byte(subAPIResp),
	}

	amsc := newTestAMSClient(t, APIResponses, 1)
	clustersInfo, err := amsc.GetClustersForOrganization("external-test-id")
	assert.Nil(t, err)

	aci, found := clustersInfo["test-external-cluster-id"]
	assert.True(t, found)
	assert.Equal(t, "test-external-cluster-id", aci.ID)
	assert.Equal(t, "test-display-name", aci.DisplayName)
	assert.Equal(t, "archived", aci.Status)
}

func TestGetSingleClusterInfoForOrganization(t *testing.T) {
	APIResponses := map[string][]byte{
		"/token":                              []byte(tokenAPIResp),
		"/api/accounts_mgmt/v1/organizations": []byte(orgAPIResp),
		"/api/accounts_mgmt/v1/subscriptions": []byte(subAPIResp),
	}

	amsc := newTestAMSClient(t, APIResponses, 1)
	clusterInfo, err := amsc.GetSingleClusterInfoForOrganization("external-test-id", "test-cluster-id")
	assert.Nil(t, err)

	aci := clusterInfo
	assert.Equal(t, "test-external-cluster-id", aci.ID)
	assert.Equal(t, "test-display-name", aci.DisplayName)
	assert.Equal(t, "archived", aci.Status)
}

func TestGetInternalOrgIDFromExternal(t *testing.T) {
	APIResponses := map[string][]byte{
		"/token":                              []byte(tokenAPIResp),
		"/api/accounts_mgmt/v1/organizations": []byte(orgAPIResp),
	}

	amsc := newTestAMSClient(t, APIResponses, 1)
	externalID, err := amsc.GetInternalOrgIDFromExternal("external-test-id")
	assert.Nil(t, err)
	assert.Equal(t, "internal-test-id", externalID)
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

	os.Exit(m.Run())
}
