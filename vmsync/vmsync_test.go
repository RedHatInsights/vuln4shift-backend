package vmsync

import (
	"app/base/models"
	"app/base/utils"
	"app/test"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const VmaasRespMock = `{
    "cve_list": {
        "CVE-2015-5320": {
            "redhat_url": "https://access.redhat.com/security/cve/cve-2015-5320",
            "secondary_url": "https://wiki.jenkins-ci.org/display/SECURITY/Jenkins+Security+Advisory+2015-11-11",
            "synopsis": "CVE-2015-5320",
            "impact": "Important",
            "public_date": "2015-11-11T00:00:00+00:00",
            "modified_date": "2022-09-20T08:12:58+00:00",
            "cwe_list": [],
            "cvss3_score": "",
            "cvss3_metrics": "",
            "cvss2_score": "6.8",
            "cvss2_metrics": "AV:N/AC:M/Au:N/C:P/I:P/A:P",
            "description": "Jenkins before 1.638 and LTS before 1.625.2 do not properly verify the shared secret",
            "package_list": [
                "php-debuginfo-5.3.3-46.el6_7.1.x86_64"
            ],
            "source_package_list": [
                "openshift-origin-cartridge-python-1.34.2.1-1.el6op.src"
            ],
            "errata_list": [
                "RHSA-2016:22381"
            ]
        }
    },
    "page": 1,
    "page_size": 1,
    "pages": 1,
    "last_change": "2022-09-20T08:39:52.919848+00:00"
}`

const VmaasCvesCnt = 1

// Test sync with preloading CVEs from DB like on the startup.
func TestSyncCveMetadataPreload(t *testing.T) {
	beforeCnt := len(test.GetAllCves(t))
	assert.Nil(t, prepareDbCvesMap())

	// CVEs from VMaaS response must be inserted into DB during this sync.
	// CVEs from test DB are not pruned because they are associated with an image.
	syncCveMetadata()
	assert.Equal(t, beforeCnt+VmaasCvesCnt, len(test.GetAllCves(t)))
}

func TestGetApiCves(t *testing.T) {
	var apiResp APICveResponse
	assert.Nil(t, json.Unmarshal([]byte(VmaasRespMock), &apiResp))
	assert.Equal(t, 1, len(apiResp.CveList))

	cveNames, cves, err := getAPICves()
	assert.Nil(t, err)

	i := 0
	for cveName, cve := range apiResp.CveList {
		// True if mock resp is sorted.
		assert.Equal(t, cveName, cveNames[i])
		i++
		assert.True(t, reflect.DeepEqual(cve, cves[cveName]))
	}
}

func TestGetApiCvesFailRequest(t *testing.T) {
	errMsg := "expected error"
	httpClient = test.NewAPIMock("Bad Request", 400, nil, errors.New(errMsg))
	_, _, err := getAPICves()
	assert.Equal(t, errMsg, err.Error())
}

func TestPruneCves(t *testing.T) {
	// CVEs unassociated with any cluster and absent in VMaaS response.
	toPruneCves := []models.Cve{
		{
			Name:        "CVE-2018-11108",
			Description: "unknown",
			Severity:    "NotSet",
		},
		{
			Name:        "VE-2018-11108",
			Description: "unknown",
			Severity:    "NotSet",
		},
	}
	id := test.InsertCve(t, toPruneCves[0])
	id2 := test.InsertCve(t, toPruneCves[1])

	beforeSyncCnt := len(test.GetAllCves(t))
	nonAffectingCnt := len(test.GetNonAffectingCves(t))

	assert.Nil(t, prepareDbCvesMap())
	assert.Nil(t, pruneCves())

	afterSyncCnt := len(test.GetAllCves(t))
	assert.True(t, beforeSyncCnt-afterSyncCnt == nonAffectingCnt)

	assert.Nil(t, test.GetCveByID(t, id))
	assert.Nil(t, test.GetCveByID(t, id2))
}

func TestSyncCves(t *testing.T) {
	now := time.Now()

	toUpdateCve := test.GetCveByID(t, 22)
	assert.NotNil(t, toUpdateCve)

	toSyncCves := []models.Cve{
		{
			ID:          0,
			Name:        "sync-cve-test-1",
			Description: "sync-cve-description",
			Severity:    "NotSet",
		},
		{
			ID:          1,
			Name:        "sync-cve-test-2",
			Description: "sync-cve-description",
			Severity:    "NotSet",
		},
		// Assuming update is taking place for this ID.
		{
			ID:           toUpdateCve.ID,
			Name:         toUpdateCve.Name,
			Description:  "sync-cve-description-update",
			Severity:     "NotSet",
			ModifiedDate: &now,
		},
	}

	beforeCves := test.GetAllCves(t)

	assert.Nil(t, syncCves(toSyncCves))

	cveAfterUpdate := test.GetCveByID(t, toUpdateCve.ID)
	assert.Equal(t, now.Format(time.RFC3339), (*cveAfterUpdate.ModifiedDate).Format(time.RFC3339))

	afterCves := test.GetAllCves(t)
	insertedCnt := len(afterCves) - len(beforeCves)
	toInsertCnt := len(toSyncCves) - 1
	assert.Equal(t, toInsertCnt, insertedCnt)

	// Restore updated CVE
	test.UpsertCves(t, []models.Cve{*toUpdateCve})

	// Remove inserted CVEs
	test.DeleteCvesByID(t, toSyncCves[0].ID, toSyncCves[1].ID)
}

func TestGetMetricsPusher(t *testing.T) {
	srv := test.GetMetricsServer(t, "PUT", "vmsync")
	defer srv.Close()

	oldPrometheusGateway := utils.Cfg.PrometheusPushGateway
	defer func() { utils.Cfg.PrometheusPushGateway = oldPrometheusGateway }()
	utils.Cfg.PrometheusPushGateway = srv.URL

	pusher := GetMetricsPusher()
	assert.Nil(t, pusher.Push())
}

func TestMain(m *testing.M) {
	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		panic(err)
	}

	test.DB = db
	DB = test.DB

	BatchSize = 5000
	PageSize = 5000

	// Mock HTTP client used to call VMaaS
	httpClient = test.NewAPIMock("OK", 200, []byte(VmaasRespMock), nil)

	err = test.ResetDB()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
