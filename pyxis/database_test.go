package pyxis

import (
	"app/base/models"
	"app/test"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbConfigure(t *testing.T) {
	defer func() { DB = test.DB }()

	assert.Nil(t, dbConfigure())
}

func TestPrepareDbRepositories(t *testing.T) {
	expectedRepos := test.GetAllRepos(t)

	assert.Nil(t, prepareDbRepositories())
	for _, expectedRepo := range expectedRepos {
		actualRepo, found := dbRepoMap[fmt.Sprintf("%s/%s", expectedRepo.Registry, expectedRepo.Repository)]
		assert.True(t, found)
		assert.Equal(t, expectedRepo, actualRepo)

		actualRepoByID, found := dbPyxisIDRepoMap[expectedRepo.PyxisID]
		assert.True(t, found)
		assert.Equal(t, expectedRepo, actualRepoByID)
	}
}

func TestPrepareDbArchs(t *testing.T) {
	expectedArchs := map[string]bool{
		"amd64": true,
	}

	assert.Nil(t, prepareDbArchs())
	for arch := range dbArchMap {
		assert.True(t, expectedArchs[arch])
	}
}

func TestPrepareDbImages(t *testing.T) {
	expectedImages := test.GetAllImages(t)

	assert.Nil(t, prepareDbImages())
	for _, expectedImg := range expectedImages {
		actualImg, found := dbImageMap[expectedImg.PyxisID]
		assert.True(t, found)
		assert.Equal(t, expectedImg, actualImg)
	}
}

func TestPrepareMaps(t *testing.T) {
	assert.Nil(t, prepareMaps())
}

func TestEmptyPendingCache(t *testing.T) {
	emptyPendingCache()
	assert.Equal(t, 0, len(dbArchMapPending))
	assert.Equal(t, 0, len(dbImageMapPending))
	assert.Equal(t, 0, len(dbCveMapPending))
}

func TestFlushPendingCache(t *testing.T) {
	assert.Nil(t, prepareDbArchs())
	assert.Nil(t, prepareDbImages())
	assert.Nil(t, prepareDbCves())

	expectedArchs := dbArchMap
	expectedImages := dbImageMap
	expectedCves := dbCveMap

	dbArchMapPending = expectedArchs
	dbImageMapPending = expectedImages
	dbCveMapPending = expectedCves

	dbArchMap = map[string]models.Arch{}
	dbImageMap = map[string]models.Image{}
	dbCveMap = map[string]models.Cve{}

	flushPendingCache()

	assert.Equal(t, expectedArchs, dbArchMap)
	assert.Equal(t, expectedImages, dbImageMap)
	assert.Equal(t, expectedCves, dbCveMap)
}

func TestGetDbImageCves(t *testing.T) {
	for _, subjectImage := range test.GetAllImages(t) {
		expectedImageCves := test.GetAllImageCves(t, subjectImage.ID)

		actualCves, err := getDbImageCves(subjectImage.ID)
		assert.Nil(t, err)

		for _, expectedCve := range expectedImageCves {
			actualCve, found := actualCves[expectedCve.CveID]
			assert.True(t, found)
			assert.Equal(t, expectedCve, actualCve)
		}
	}
}
