package pyxis

import (
	"app/base/models"
	"app/base/utils"
	"app/test"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIRepositories(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposResp), nil)

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
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIRepoImagesResp), nil)

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

	assert.Equal(t, expectedDockerDigest, image.DockerImageDigest)
	assert.Equal(t, expectedArch, image.Architecture)
	assert.Equal(t, expectedModified, image.ModifiedDate.UTC().String())
}

func TestGetAPIRepoImagesBadRequest(t *testing.T) {
	httpClient = test.NewAPIMock("Bad Request", 400, nil, errors.New("bad request"))

	registry := "registry.access.redhat.com"
	repo := "rhel7.1"

	_, err := getAPIRepoImages(registry, repo)
	assert.Equal(t, "bad request", err.Error())
}

func TestGetAPIRepoImagesNoRepos(t *testing.T) {
	repoImages := APIRepoImagesResponse{}
	assert.Nil(t, json.Unmarshal([]byte(test.PyxisAPIRepoImagesResp), &repoImages))
	repoImages.Data[0].Repositories = []APIImageRepoDetail{}

	repoImagesBytes, err := json.Marshal(repoImages)
	assert.Nil(t, err)

	httpClient = test.NewAPIMock("OK", 200, repoImagesBytes, nil)

	registry := "registry.access.redhat.com"
	repo := "rhel7.1"

	_, err = getAPIRepoImages(registry, repo)
	assert.Equal(t, fmt.Sprintf("Empty repositories field for Image Pyxis ID: %s", repoImages.Data[0].PyxisID), err.Error())
}

func TestGetAPIRepoImagesNoDigest(t *testing.T) {
	repoImages := APIRepoImagesResponse{}
	assert.Nil(t, json.Unmarshal([]byte(test.PyxisAPIRepoImagesRespNoDigest), &repoImages))
	repoImages.Data[0].DockerImageDigest = ""

	repoImagesBytes, err := json.Marshal(repoImages)
	assert.Nil(t, err)

	httpClient = test.NewAPIMock("OK", 200, repoImagesBytes, nil)

	registry := "registry.access.redhat.com"
	repo := "rhel7.1"

	_, err = getAPIRepoImages(registry, repo)
	assert.Equal(t, fmt.Sprintf("Empty manifest_list_digest, manifest_schema2_digest and docker_image_digest fields for Image Pyxis ID: %s", repoImages.Data[0].PyxisID), err.Error())
}

func TestGetAPIImageCves(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIImageCvesResp), nil)

	imageID := "57ea8d0d9c624c035f96f45e"
	expectedCves := []string{"CVE-2016-2180"}

	resp, err := getAPIImageCves(imageID)
	assert.Nil(t, err)
	for i, cve := range expectedCves {
		assert.Equal(t, cve, resp[i])
	}
}

func TestGetAPIImageCvesBadRequest(t *testing.T) {
	httpClient = test.NewAPIMock("Bad Request", 400, nil, errors.New("bad request"))

	imageID := "57ea8d0d9c624c035f96f45e"

	_, err := getAPIImageCves(imageID)
	assert.Equal(t, "bad request", err.Error())
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
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIImageCvesResp), nil)
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

	newImageID := test.GetImageByPyxisID(t, pyxisID).ID
	test.DeleteImageCves(t, newImageID)
	test.DeleteImage(t, pyxisID)
}

func TestSyncRepoNew(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposRespJBoss), nil)
	assert.Nil(t, prepareDbCves())
	test.DeleteRepo(t, "jboss-fuse-6", "registry.access.redhat.com")

	repoToSync := models.Repository{
		PyxisID:      "test-pyxis-jboss",
		ModifiedDate: time.Time{},
		Registry:     "registry.access.redhat.com",
		Repository:   "jboss-fuse-8",
	}
	assert.Nil(t, syncRepo(repoToSync))

	syncedRepo := test.GetRepo(t, repoToSync.Registry, repoToSync.Repository)
	assert.Equal(t, repoToSync.Repository, syncedRepo.Repository)
	assert.Equal(t, repoToSync.ModifiedDate.UTC(), syncedRepo.ModifiedDate.UTC())
	assert.Equal(t, repoToSync.PyxisID, syncedRepo.PyxisID)
	assert.Equal(t, repoToSync.Registry, syncedRepo.Registry)
}

// Should update existing repo, insert new Image and new CVE associations.
func TestSyncRepo(t *testing.T) {
	var repoImagesMockResp APIRepoImagesResponse
	assert.Nil(t, json.Unmarshal([]byte(test.PyxisAPIRepoImagesNewResp), &repoImagesMockResp))
	expectedAPIImages := repoImagesMockResp.Data

	var imageCvesMockResp APIImageCvesResponse
	assert.Nil(t, json.Unmarshal([]byte(test.PyxisAPIImageCvesResp), &imageCvesMockResp))

	var expectedCveNames []string
	for _, APICve := range imageCvesMockResp.Data {
		expectedCveNames = append(expectedCveNames, APICve.Cve)
	}

	assert.Nil(t, prepareDbCves())
	assert.Nil(t, prepareMaps())

	registry := "registry.access.redhat.com"
	repository := "jboss-fuse-6"
	pyxisID := "57ea8d0d9c624c488f96f45e"

	repoToSync := models.Repository{
		ID:           7357,
		PyxisID:      pyxisID,
		ModifiedDate: time.Time{},
		Registry:     registry,
		Repository:   repository,
	}

	test.CreateRepo(t, repoToSync)

	mockResponses := map[string][]byte{
		fmt.Sprintf("/repositories/registry/%s/repository/%s/images", registry, repository): []byte(test.PyxisAPIRepoImagesNewResp),
		fmt.Sprintf("/images/id/%s/vulnerabilities", pyxisID):                               []byte(test.PyxisAPIImageCvesResp),
	}
	httpClient = test.NewAPIMockMultiEndpoint("OK", 200, mockResponses, nil)

	assert.Nil(t, syncRepo(repoToSync))

	actualRepos := test.GetAllRepos(t)

	// Check updated repos
	for _, actualRepo := range actualRepos {
		for _, image := range expectedAPIImages {
			if actualRepo.PyxisID == image.PyxisID {
				assert.Equal(t, actualRepo.Repository, repository)
				assert.Equal(t, actualRepo.Registry, registry)
				assert.True(t, image.ModifiedDate.After(actualRepo.ModifiedDate))
			}
		}
	}

	expectedCves := test.GetCvesByName(t, expectedCveNames...)
	newImageID := test.GetImageByPyxisID(t, pyxisID).ID
	actualCves := test.GetAllImageCves(t, newImageID)

	// Check new associated CVEs.
	for i, actualCve := range actualCves {
		assert.Equal(t, expectedCves[i].ID, actualCve.CveID)
	}

	test.DeleteRepoImageByImageID(t, newImageID)
	test.DeleteImageCves(t, newImageID)
	test.DeleteImage(t, pyxisID)
	test.DeleteRepo(t, registry, repository)
}

// Should create new repo, new Image and new CVE associations.
func TestSyncRepoNewImage(t *testing.T) {
	var repoImagesMockResp APIRepoImagesResponse
	assert.Nil(t, json.Unmarshal([]byte(test.PyxisAPIRepoImagesNewResp), &repoImagesMockResp))
	expectedAPIImages := repoImagesMockResp.Data

	var imageCvesMockResp APIImageCvesResponse
	assert.Nil(t, json.Unmarshal([]byte(test.PyxisAPIImageCvesResp), &imageCvesMockResp))

	var expectedCveNames []string
	for _, APICve := range imageCvesMockResp.Data {
		expectedCveNames = append(expectedCveNames, APICve.Cve)
	}

	assert.Nil(t, prepareDbCves())
	assert.Nil(t, prepareMaps())

	registry := "registry.access.redhat.com"
	repository := "jboss-fuse-6"
	newPyxisID := "57ea8d0d9c624c488f96f45e"

	repoToSync := models.Repository{
		PyxisID:      newPyxisID,
		ModifiedDate: time.Time{},
		Registry:     registry,
		Repository:   repository,
	}

	delete(dbImageMap, newPyxisID)

	mockResponses := map[string][]byte{
		fmt.Sprintf("/repositories/registry/%s/repository/%s/images", registry, repository): []byte(test.PyxisAPIRepoImagesNewResp),
		fmt.Sprintf("/images/id/%s/vulnerabilities", newPyxisID):                            []byte(test.PyxisAPIImageCvesResp),
	}
	httpClient = test.NewAPIMockMultiEndpoint("OK", 200, mockResponses, nil)

	assert.Nil(t, syncRepo(repoToSync))

	actualRepos := test.GetAllRepos(t)

	// Check updated repos
	for _, actualRepo := range actualRepos {
		for _, image := range expectedAPIImages {
			if actualRepo.PyxisID == image.PyxisID {
				assert.Equal(t, actualRepo.Repository, repository)
				assert.Equal(t, actualRepo.Registry, registry)
				assert.True(t, image.ModifiedDate.After(actualRepo.ModifiedDate))
			}
		}
	}

	expectedCves := test.GetCvesByName(t, expectedCveNames...)
	newImageID := test.GetImageByPyxisID(t, newPyxisID).ID
	actualCves := test.GetAllImageCves(t, newImageID)

	// Check new associated CVEs.
	for i, actualCve := range actualCves {
		assert.Equal(t, expectedCves[i].ID, actualCve.CveID)
	}

	test.DeleteRepoImageByImageID(t, newImageID)
	test.DeleteImageCves(t, newImageID)
	test.DeleteImage(t, newPyxisID)
	test.DeleteRepo(t, registry, repository)
}

// Should delete repository images not returned by an API.
func TestSyncReposDelete(t *testing.T) {
	assert.Nil(t, prepareDbCves())
	assert.Nil(t, prepareMaps())

	registry := "registry.access.redhat.com"
	repository := "jboss-fuse-6"
	newPyxisID := "57ea8d0d9c624c488f96f45e"

	repoToSync := models.Repository{
		ID:           7357,
		PyxisID:      newPyxisID,
		ModifiedDate: time.Time{},
		Registry:     registry,
		Repository:   repository,
	}

	test.CreateRepo(t, repoToSync)

	imageToRemove := models.Image{
		ID:           7357,
		PyxisID:      "dummy-pyxis-id",
		ModifiedDate: time.Now(),
		ArchID:       1,
	}
	test.CreateImage(t, imageToRemove)
	test.CreateRepoImage(t, models.RepositoryImage{
		RepositoryID: repoToSync.ID,
		ImageID:      imageToRemove.ID,
	})

	mockResponses := map[string][]byte{
		fmt.Sprintf("/repositories/registry/%s/repository/%s/images", registry, repository): []byte(test.PyxisAPIRepoImagesNewResp),
		fmt.Sprintf("/images/id/%s/vulnerabilities", newPyxisID):                            []byte(test.PyxisAPIImageCvesResp),
	}
	httpClient = test.NewAPIMockMultiEndpoint("OK", 200, mockResponses, nil)

	assert.Nil(t, syncRepo(repoToSync))

	actualRepoImages := test.GetRepoImages(t, repoToSync.ID)
	for _, repoImage := range actualRepoImages {
		assert.NotEqual(t, imageToRemove.ID, repoImage.ImageID)
	}

	newImageID := test.GetImageByPyxisID(t, newPyxisID).ID
	test.DeleteRepoImageByImageID(t, newImageID)
	test.DeleteImageCves(t, newImageID)
	test.DeleteImage(t, newPyxisID)
	test.DeleteRepo(t, registry, repository)
}

func TestSyncRepos(t *testing.T) {
	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposRespSync), nil)
	profile = testProfile
	parseProfiles()
	assert.Nil(t, prepareMaps())

	delete(dbRepoMap, "registry.access.redhat.com/rhel6")
	delete(dbPyxisIDRepoMap, "registry.access.redhat.com/rhel7/sadc")

	syncRepos()
}

func TestStart(t *testing.T) {
	srv := test.GetMetricsServer(t, "POST", "pyxis")
	defer srv.Close()

	oldPrometheusGateway := utils.Cfg.PrometheusPushGateway
	defer func() { utils.Cfg.PrometheusPushGateway = oldPrometheusGateway }()
	utils.Cfg.PrometheusPushGateway = srv.URL

	httpClient = test.NewAPIMock("OK", 200, []byte(test.PyxisAPIReposRespJBoss), nil)
	profile = testProfile

	Start()
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
