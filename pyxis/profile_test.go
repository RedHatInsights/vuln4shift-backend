package pyxis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testProfile  = "testing"
	testRegistry = "registry.access.redhat.com"
	testRepos    = map[string]struct{}{
		"ubi7/ubi":              {},
		"ubi7/ubi-minimal":      {},
		"ubi8/ubi":              {},
		"ubi8/ubi-minimal":      {},
		"ubi8/ubi-micro":        {},
		"ubi9-beta/ubi":         {},
		"ubi9-beta/ubi-minimal": {},
		"ubi9-beta/ubi-micro":   {},
		"rhel7.1":               {},
		"rhel7/sadc":            {},
		"rhel6":                 {},
	}
)

func TestParseProfiles(t *testing.T) {
	profile = testProfile
	parseProfiles()

	actualRegistry := profileMap[testProfile]
	assert.NotNil(t, actualRegistry)

	actualRepos := actualRegistry[testRegistry]
	assert.Equal(t, len(testRepos), len(actualRepos))
	assert.Equal(t, testRepos, actualRepos)
}

func TestRepositoryInProfile(t *testing.T) {
	// Load testing profileMap
	profile = testProfile
	parseProfiles()

	for repo := range testRepos {
		assert.True(t, repositoryInProfile(testRegistry, repo))
	}
}
