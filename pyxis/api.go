package pyxis

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"app/base/api"
	"app/base/utils"
)

var (
	PyxisBaseURL       = utils.Cfg.PyxisBaseURL
	PyxisReposURL      = fmt.Sprintf("%s/repositories", PyxisBaseURL)
	PyxisRepoImagesURL = fmt.Sprintf("%s/repositories/registry/%%s/repository/%%s/images", PyxisBaseURL)
	PyxisImageCvesURL  = fmt.Sprintf("%s/images/id/%%s/vulnerabilities", PyxisBaseURL)
	PageSize           = 500
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

type APIImage struct {
	PyxisID      string    `json:"_id"`
	ModifiedDate time.Time `json:"last_update_date"`
	Digest       string    `json:"docker_image_id"`
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

	client := &api.Client{HTTPClient: &http.Client{}}
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
			repoMap[repo.PyxisID] = repo
		}

		totalPages = getTotalPages(pyxisResponse.Total)
		logger.Infof("Fetched Pyxis repository list: repos=%d, page=%d/%d", len(repoMap), page+1, totalPages)
	}

	return repoMap, nil
}

func getAPIRepoImages(registry, repository string) (map[string]APIImage, error) {
	imageMap := make(map[string]APIImage)

	client := &api.Client{HTTPClient: &http.Client{}}
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
			if len(image.Digest) == 0 {
				err := fmt.Errorf("Empty digest field for Image Pyxis ID: %s", image.PyxisID)
				return imageMap, err // Break here, do not sync repo if at least one image is faulty
			}
			// De-duplicate on the digest field, not unique in Pyxis DB, but we need only one record
			// all records with same digest in Pyxis should have same CVE list
			imageMap[image.Digest] = image
		}

		totalPages = getTotalPages(pyxisResponse.Total)
		logger.Infof("Fetched Pyxis repository images: repo=%s/%s, images=%d, page=%d/%d", registry, repository, len(imageMap), page+1, totalPages)
	}

	return imageMap, nil
}

func getAPIImageCves(imagePyxisID string) (map[string]struct{}, error) {
	cveMap := make(map[string]struct{})

	client := &api.Client{HTTPClient: &http.Client{}}
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
			return cveMap, err
		}

		for _, cve := range pyxisResponse.Data {
			cveMap[cve.Cve] = struct{}{} // We need only key - CVE name
		}

		totalPages = getTotalPages(pyxisResponse.Total)
		logger.Infof("Fetched Pyxis image CVEs: image=%s, cves=%d, page=%d/%d", imagePyxisID, len(cveMap), page+1, totalPages)
	}

	return cveMap, nil
}
