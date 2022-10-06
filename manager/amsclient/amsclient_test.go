package amsclient

import (
	"app/base/ams"
	"app/base/models"
	"app/base/utils"
	"app/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBFetchClusterDetails(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		expectedClusters := test.GetAccountClusters(t, account.ID)

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
		amsMock := &test.AMSClientMock{
			ClustersResponse: amsClusters,
		}

		uuids, statuses, versions, providers, err := DBFetchClusterDetails(test.DB, amsMock, account.ID, account.OrgID, true, nil)
		assert.Nil(t, err)

		for i, uuid := range uuids {
			assert.Equal(t, uuid, expectedClusters[i].UUID.String())
		}

		if len(uuids) > 0 {
			test.CheckClusterDetails(t, expectedProviders, expectedStatuses, expectedVersions, providers, statuses, versions)
		}
	}
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
