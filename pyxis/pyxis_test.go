package pyxis

import (
	"app/base/models"
	"app/base/utils"
	"app/test"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIRepositories(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposResp))

	resp, err := getAPIRepositories()
	assert.Nil(t, err)

	expectedRegistry := "registry.access.redhat.com"
	expectedRepo := "rhel7.1"
	expectedID := "57ea8cd79c624c035f96f327"
	expectedModified := "2022-05-12 11:30:50.822 +0000 UTC"

	repo, exists := resp[fmt.Sprintf("%s/%s", expectedRegistry, expectedRepo)]
	assert.True(t, exists)
	assert.Equal(t, expectedRepo, repo.Repository)
	assert.Equal(t, expectedID, repo.PyxisID)
	assert.Equal(t, expectedRegistry, repo.Registry)
	assert.Equal(t, expectedModified, repo.ModifiedDate.UTC().String())
}

func TestGetAPIRepoImages(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIRepoImagesResp))

	registry := "registry.access.redhat.com"
	repo := "rhel7.1"

	resp, err := getAPIRepoImages(registry, repo)
	assert.Nil(t, err)

	expectedID := "57ea8d0d9c624c035f96f45e"
	expectedDockerDigest := "temp:sha256:3817ddfacc32be3501dce396efcbf864ec68c3d9794a38d0c959377fca65e881"
	expectedArch := "amd64"
	expectedModified := "2022-10-07 01:51:21.689 +0000 UTC"

	image, exists := resp[fmt.Sprintf(expectedID)]
	assert.True(t, exists)
	//assert.Equal(t, nil, image.Repositories)
	assert.Equal(t, expectedDockerDigest, image.DockerImageDigest)
	assert.Equal(t, expectedArch, image.Architecture)
	assert.Equal(t, expectedModified, image.ModifiedDate.UTC().String())
}

func TestGetAPIImageCves(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIImageCvesResp))

	imageID := "57ea8d0d9c624c035f96f45e"
	expectedCves := []string{"CVE-2016-2180"}

	resp, err := getAPIImageCves(imageID)
	assert.Nil(t, err)
	for i, cve := range expectedCves {
		assert.Equal(t, cve, resp[i])
	}
}

func TestRegisterMissingCves(t *testing.T) {
	assert.Nil(t, prepareDbCves())

	cves := []string{"CVE-2016-2180", "CVE-2016-7141", "CVE-2016-3075"}
	assert.Nil(t, registerMissingCves(test.DB, cves))

	for _, cve := range cves {
		_, found := dbCveMapPending[cve]
		assert.True(t, found)
	}
}

func TestSyncImage(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIImageCvesResp))
	assert.Nil(t, prepareDbCves())

	var ID int64 = 7357
	pyxisID := "test-pyxis-id"
	ModifiedDate := time.Now()
	DockerImageDigest := "7357"
	var ArchID int64 = 1
	img := models.Image{ID: ID, PyxisID: pyxisID, ModifiedDate: ModifiedDate, DockerImageDigest: &DockerImageDigest, ArchID: ArchID}

	test.DeleteImage(t, pyxisID)

	assert.Nil(t, syncImage(test.DB, img))

	insertedImg := test.GetImageByID(t, ID)
	assert.Equal(t, pyxisID, insertedImg.PyxisID)
	assert.Equal(t, ModifiedDate.Round(time.Second), insertedImg.ModifiedDate.Round(time.Second))
	assert.Equal(t, DockerImageDigest, *insertedImg.DockerImageDigest)
	assert.Equal(t, ArchID, insertedImg.ArchID)
}

func TestSyncRepoNew(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposRespJBoss))
	assert.Nil(t, prepareDbCves())
	test.DeleteRepo(t, "jboss-fuse-67", "registry.access.redhat.com")

	repoToSync := models.Repository{
		PyxisID:      "test-pyxis-jboss",
		ModifiedDate: time.Time{},
		Registry:     "registry.access.redhat.com",
		Repository:   "jboss-fuse-6",
	}
	assert.Nil(t, syncRepo(repoToSync))

	syncedRepo := test.GetRepo(t, repoToSync.Repository)
	assert.Equal(t, repoToSync.Repository, syncedRepo.Repository)
	assert.Equal(t, repoToSync.ModifiedDate.UTC(), syncedRepo.ModifiedDate.UTC())
	assert.Equal(t, repoToSync.PyxisID, syncedRepo.PyxisID)
	assert.Equal(t, repoToSync.Registry, syncedRepo.Registry)
}

func TestSyncRepo(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposRespJBoss))
	assert.Nil(t, prepareDbCves())

	registry := "registry.access.redhat.com"
	repository := "jboss-fuse-6"
	repoToSync := models.Repository{
		ID:           20,
		PyxisID:      "test-pyxis-jboss",
		ModifiedDate: time.Time{},
		Registry:     registry,
		Repository:   repository,
	}
	test.DeleteRepo(t, registry, repository)
	test.CreateRepo(t, repoToSync)

	assert.Nil(t, syncRepo(repoToSync))
}

func TestMain(m *testing.M) {
	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		panic(err)
	}

	test.DB = db
	DB = test.DB
	err = test.ResetDB()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
