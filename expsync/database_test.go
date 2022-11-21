package expsync

import (
	"app/base/models"
	"app/test"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exploitsTestData = `{
    "CVE-2022-0001": [
        {
            "date": "2022-01-01",
            "reference": "N/A",
            "source": "CISA"
        },
        {
            "date": "2022-01-02",
            "reference": "N/A",
            "source": "CISA"
        }
    ],
    "CVE-2022-0002": [
        {
            "date": "2022-01-03",
            "reference": "N/A",
            "source": "CISA"
        }
    ],
	"CVE-2022-9998": [
		{
			"date": "2022-01-03",
			"reference": "N/A",
			"source": "CISA"
		}
	],
	"CVE-2022-9999": [
		{}
	]
}`

func TestUpdateExploitsMetadata(t *testing.T) {
	var exploitsData map[CVE][]ExploitMetadata
	assert.Nil(t, json.Unmarshal([]byte(exploitsTestData), &exploitsData))
	const expectedSyncCnt = int64(2)

	expectedCves := make([]string, 0, len(exploitsData))
	for cve := range exploitsData {
		expectedCves = append(expectedCves, string(cve))
	}

	updatedCnt, err := updateExploitsMetadata(DB, exploitsData)
	assert.Nil(t, err)
	assert.Equal(t, expectedSyncCnt, updatedCnt)

	actualCves := test.GetCvesByName(t, expectedCves...)

	for _, actualCve := range actualCves {
		expectedMetadata, err := json.Marshal(exploitsData[CVE(actualCve.Name)])
		assert.Nil(t, err)

		actualMetadata := new(bytes.Buffer)
		assert.Nil(t, json.Compact(actualMetadata, actualCve.ExploitData))

		assert.Equal(t, expectedMetadata, actualMetadata.Bytes())
	}

	// Assert CVEs not existing in the database were not synced
	notExpectedCves := []string{"CVE-2022-9998", "CVE-2022-9999"}
	dbCves := test.GetAllCves(t)
	for _, dbCve := range dbCves {
		for _, notExpectedCve := range notExpectedCves {
			assert.NotEqual(t, dbCve.Name, notExpectedCve)
		}
	}
}

func TestGetCvesWithExploitMetadata(t *testing.T) {
	cvesWithExploit := 0

	allCves := test.GetAllCves(t)
	for _, cve := range allCves {
		if cve.ExploitData != nil {
			cvesWithExploit++
		}
	}

	actualCves, err := getCvesWithExploitMetadata(DB)
	assert.Nil(t, err)

	assert.Equal(t, cvesWithExploit, len(actualCves))
}

func TestRemoveExploitData(t *testing.T) {
	cvesBefore := test.GetAllCves(t)

	var subjects []models.Cve
	for _, cve := range cvesBefore {
		if cve.ExploitData != nil {
			subjects = append(subjects, cve)
		}
	}

	removedCnt, err := removeExploitData(DB, subjects)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(subjects)), removedCnt)

	cvesAfter := test.GetAllCves(t)

	for _, cve := range cvesAfter {
		assert.Nil(t, cve.ExploitData)
	}
}
