package expsync

import (
	"app/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncExploits(t *testing.T) {
	assert.Nil(t, syncExploits())

	cves := test.GetAllCves(t)
	for _, cve := range cves {
		if _, found := expectedAPICVEs[cve.Name]; found {
			assert.NotEqual(t, "", cve.ExploitData)
		}
	}
}
