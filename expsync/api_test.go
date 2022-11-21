package expsync

import (
	"app/base/api"
	"app/base/models"
	"app/base/utils"
	"app/test"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	exploitsRespMock = `{
		"type": "file",
		"encoding": "base64",
		"size": 5362,
		"name": "exploits.json",
		"path": "exploits.json",
		"content": "ewoJCSJDVkUtMjAyMi0wMDAxIjogWwoJCQl7CgkJCQkiZGF0ZSI6ICIyMDIyLTAzLTA3IiwKCQkJCSJyZWZlcmVuY2UiOiAiTi9BIiwKCQkJCSJzb3VyY2UiOiAiQ0lTQSIKCQkJfSwKCQkJewoJCQkJImRhdGUiOiAiMjAyMi0wMy0yNiIsCgkJCQkicmVmZXJlbmNlIjogIk4vQSIsCgkJCQkic291cmNlIjogIkNJU0EiCgkJCX0KCQldLAoJCSJDVkUtMjAyMi0wMDAyIjogWwoJCQl7CgkJCQkiZGF0ZSI6ICIyMDIyLTA4LTE4IiwKCQkJCSJyZWZlcmVuY2UiOiAiTi9BIiwKCQkJCSJzb3VyY2UiOiAiQ0lTQSIKCQkJfQoJCV0sCgkJIkNWRS0yMDIyLTAwMDMiOiBbCgkJCXsKCQkJCSJkYXRlIjogIjIwMjItMDktMDgiLAoJCQkJInJlZmVyZW5jZSI6ICJOL0EiLAoJCQkJInNvdXJjZSI6ICJDSVNBIgoJCQl9CgkJXSwKCQkiQ1ZFLTIwMjItMDAwNCI6IFsKCQkJewoJCQkJImRhdGUiOiAiMjAyMi0wOC0xOCIsCgkJCQkicmVmZXJlbmNlIjogIk4vQSIsCgkJCQkic291cmNlIjogIkNJU0EiCgkJCX0KCQldLAoJCSJDVkUtMjAyMi0wMDA1IjogWwoJCQl7CgkJCQkiZGF0ZSI6ICIyMDIyLTEwLTI4IiwKCQkJCSJyZWZlcmVuY2UiOiAiTi9BIiwKCQkJCSJzb3VyY2UiOiAiQ0lTQSIKCQkJfQoJCV0KCX0=\n",
		"sha": "3d21ec53a331a6f037a91c368710b99387d012c1"
	}`
	expectedRef    = "N/A"
	expectedSource = "CISA"
)

var expectedAPICVEs = map[string][]string{
	"CVE-2022-0001": {"2022-03-07", "2022-03-26"},
	"CVE-2022-0002": {"2022-08-18"},
	"CVE-2022-0003": {"2022-09-08"},
	"CVE-2022-0004": {"2022-08-18"},
	"CVE-2022-0005": {"2022-10-28"},
}

func TestGetAPIExploits(t *testing.T) {
	var r api.GithubRepoAPIResponse
	assert.Nil(t, json.Unmarshal([]byte(exploitsRespMock), &r))

	exps, err := getAPIExploits()
	assert.Nil(t, err)

	for cve, actualExps := range exps {
		expectedDates, found := expectedAPICVEs[string(cve)]
		assert.True(t, found)

		for i, expectedDate := range expectedDates {
			assert.Equal(t, expectedDate, actualExps[i].Date)
			assert.Equal(t, expectedSource, actualExps[i].Source)
			assert.Equal(t, expectedRef, actualExps[i].Reference)
		}
	}
}

func TestMain(m *testing.M) {
	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		panic(err)
	}

	test.DB = db
	DB = test.DB

	// Mock HTTP client used to call VMaaS
	httpClient = test.NewAPIMock("OK", 200, []byte(exploitsRespMock), nil)

	err = test.ResetDB()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
