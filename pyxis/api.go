package pyxis

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"app/base/api"
	"app/base/utils"
)

var (
	PyxisBaseURL                      = utils.Cfg.PyxisBaseURL
	PyxisReposURL                     = fmt.Sprintf("%s/repositories", PyxisBaseURL)
	PyxisRepoImagesURL                = fmt.Sprintf("%s/repositories/registry/%%s/repository/%%s/images", PyxisBaseURL)
	PyxisImageCvesURL                 = fmt.Sprintf("%s/images/id/%%s/vulnerabilities", PyxisBaseURL)
	PageSize                          = 500
	httpClient         api.HTTPClient = &http.Client{}
)

type APIRepo struct {
	PyxisID      string    `json:"_id"`
	ModifiedDate time.Time `json:"last_update_date"`
	Registry     string    `json:"registry"`
	Repository   string    `json:"repository"`
}

type APIReposResponse struct {
	Data     []APIRepo `json:"data"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Total    int       `json:"total"`
}

type APIImageRepoTag struct {
	Name string `json:"name"`
}

type APIImageRepoDetail struct {
	ManifestListDigest    string            `json:"manifest_list_digest"`
	ManifestSchema2Digest string            `json:"manifest_schema2_digest"`
	Tags                  []APIImageRepoTag `json:"tags"`
}

type APIImage struct {
	PyxisID           string               `json:"_id"`
	ModifiedDate      time.Time            `json:"last_update_date"`
	Architecture      string               `json:"architecture"`
	DockerImageDigest string               `json:"docker_image_digest"`
	Repositories      []APIImageRepoDetail `json:"repositories"`
}

type APIRepoImagesResponse struct {
	Data     []APIImage `json:"data"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
}

type APICve struct {
	Cve string `json:"cve_id"`
}

type APIImageCvesResponse struct {
	Data     []APICve `json:"data"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Total    int      `json:"total"`
}

func getTotalPages(totalItems int) int {
	totalPages := int(math.Ceil(float64(totalItems) / float64(PageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	return totalPages
}

func getAPIRepositories() (map[string]APIRepo, error) {
	repoMap := make(map[string]APIRepo)

	client := &api.Client{HTTPClient: httpClient}
	totalPages := 9999

	// Pyxis indexes pages from 0
	for page := 0; page < totalPages; page++ {
		reposURL := fmt.Sprintf("%s?page_size=%d&page=%d", PyxisReposURL, PageSize, page)
		pyxisResponse := APIReposResponse{}
		statusCode, err := client.RetryRequest(http.MethodGet, reposURL, nil, &pyxisResponse)
		if err != nil {
			pyxisRequestError.WithLabelValues(reposURL, http.MethodGet, strconv.Itoa(statusCode)).Inc()
			logger.Warningf("Request %s %s failed: statusCode=%d, err=%s", http.MethodGet, reposURL, statusCode, err)
			return repoMap, err
		}

		for _, repo := range pyxisResponse.Data {
			if len(repo.Registry) == 0 {
				logger.Debugf("Empty registry field for Repository Pyxis ID: %s", repo.PyxisID)
				continue
			}
			if len(repo.Repository) == 0 {
				logger.Debugf("Empty repository field for Repository Pyxis ID: %s", repo.PyxisID)
				continue
			}
			repoMap[formatRepoMapKey(repo.Registry, repo.Repository)] = repo
		}

		totalPages = getTotalPages(pyxisResponse.Total)
		logger.Infof("Fetched Pyxis repository list: repos=%d, page=%d/%d", len(repoMap), page+1, totalPages)
	}

	return repoMap, nil
}

func getAPIRepoImages(registry, repository string) (map[string]APIImage, error) {
	imageMap := make(map[string]APIImage)

	client := &api.Client{HTTPClient: httpClient}
	repoImagesURL := fmt.Sprintf(PyxisRepoImagesURL, registry, repository)
	totalPages := 9999

	// Pyxis indexes pages from 0
	for page := 0; page < totalPages; page++ {
		repoImagesPageURL := fmt.Sprintf("%s?page_size=%d&page=%d", repoImagesURL, PageSize, page)
		pyxisResponse := APIRepoImagesResponse{}
		statusCode, err := client.RetryRequest(http.MethodGet, repoImagesPageURL, nil, &pyxisResponse)
		if err != nil {
			pyxisRequestError.WithLabelValues(repoImagesPageURL, http.MethodGet, strconv.Itoa(statusCode)).Inc()
			logger.Warningf("Request %s %s failed: statusCode=%d, err=%s", http.MethodGet, repoImagesPageURL, statusCode, err)
			return imageMap, err
		}

		for _, image := range pyxisResponse.Data {
			if len(image.Repositories) == 0 {
				err := fmt.Errorf("Empty repositories field for Image Pyxis ID: %s", image.PyxisID)
				return imageMap, err // Break here, do not sync repo if at least one image is faulty
			}
			if len(image.DockerImageDigest) == 0 {
				for _, detail := range image.Repositories {
					if len(detail.ManifestListDigest) == 0 && len(detail.ManifestSchema2Digest) == 0 {
						err := fmt.Errorf("Empty manifest_list_digest, manifest_schema2_digest and docker_image_digest fields for Image Pyxis ID: %s", image.PyxisID)
						return imageMap, err // Break here, do not sync repo if at least one image is faulty
					}
				}
			}
			imageMap[image.PyxisID] = image
		}

		totalPages = getTotalPages(pyxisResponse.Total)
		logger.Infof("Fetched Pyxis repository images: repo=%s/%s, images=%d, page=%d/%d", registry, repository, len(imageMap), page+1, totalPages)
	}

	return imageMap, nil
}

func getAPIImageCves(imagePyxisID string) ([]string, error) {
	cveList := []string{}
	cveMap := make(map[string]struct{})

	client := &api.Client{HTTPClient: httpClient}
	imageCvesURL := fmt.Sprintf(PyxisImageCvesURL, imagePyxisID)
	totalPages := 9999

	// Pyxis indexes pages from 0
	for page := 0; page < totalPages; page++ {
		imageCvesPageURL := fmt.Sprintf("%s?page_size=%d&page=%d", imageCvesURL, PageSize, page)
		pyxisResponse := APIImageCvesResponse{}
		statusCode, err := client.RetryRequest(http.MethodGet, imageCvesPageURL, nil, &pyxisResponse)
		if err != nil {
			pyxisRequestError.WithLabelValues(imageCvesPageURL, http.MethodGet, strconv.Itoa(statusCode)).Inc()
			logger.Warningf("Request %s %s failed: statusCode=%d, err=%s", http.MethodGet, imageCvesPageURL, statusCode, err)
			return cveList, err
		}

		for _, cve := range pyxisResponse.Data {
			cveMap[cve.Cve] = struct{}{}
		}

		totalPages = getTotalPages(pyxisResponse.Total)
		logger.Infof("Fetched Pyxis image CVEs: image=%s, cves=%d, page=%d/%d", imagePyxisID, len(cveMap), page+1, totalPages)
	}

	cveList = make([]string, 0, len(cveMap))
	for cve := range cveMap {
		if !utils.IsValidCVE(cve) {
			logger.Warnf("Invalid CVE obtained from pyxis: %s, skipping", cve)
			continue
		}
		cveList = append(cveList, cve)
	}
	sort.Strings(cveList)

	return cveList, nil
}
