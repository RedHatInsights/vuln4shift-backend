package vmsync

import (
	"app/base/api"
	"app/base/utils"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var (
	PageSize     = utils.Cfg.VmaasPageSize
	VmaasCvesURL = utils.Cfg.VmaasBaseURL
	// HTTP client used to call VMaaS
	httpClient api.HTTPClient = &http.Client{}
)

type APICveRequest struct {
	CveList          []string `json:"cve_list"`
	Page             int      `json:"page"`
	PageSize         int      `json:"page_size"`
	RhOnly           bool     `json:"rh_only"`
	ErrataAssociated bool     `json:"errata_associated"`
}

type APICveResponse struct {
	CveList          map[string]APICve `json:"cve_list"`
	Page             int               `json:"page"`
	PageSize         int               `json:"page_size"`
	Pages            int               `json:"pages"`
	RhOnly           bool              `json:"rh_only"`
	ErrataAssociated bool              `json:"errata_associated"`
}

type APICve struct {
	RedhatURL         string    `json:"redhat_url"`
	SecondaryURL      string    `json:"secondary_url"`
	Synopsis          string    `json:"synopsis"`
	Impact            string    `json:"impact"`
	PublicDate        vmaasTime `json:"public_date"`
	ModifiedDate      vmaasTime `json:"modified_date"`
	CweList           []string  `json:"cwe_list"`
	Cvss3Score        string    `json:"cvss3_score"`
	Cvss3Metrics      string    `json:"cvss3_metrics"`
	Cvss2Score        string    `json:"cvss2_score"`
	Cvss2Metrics      string    `json:"cvss2_metrics"`
	Description       string    `json:"description"`
	PackageList       []string  `json:"package_list"`
	SourcePackageList []string  `json:"source_package_list"`
	ErrataList        []string  `json:"errata_list"`
}

// vmaasTime to handle empty string for time.Time fields
type vmaasTime struct {
	*time.Time
}

func (t *vmaasTime) UnmarshalJSON(b []byte) error {
	if string(b) == `""` {
		*t = vmaasTime{Time: &time.Time{}}
		return nil
	}

	str := string(b)
	str = str[1 : len(str)-1] // Trims the surrounding double quotes

	time, err := time.Parse(time.RFC3339, str)
	t.Time = &time
	return err
}

// getAPICves request CVE list from VMaaS
func getAPICves() ([]string, map[string]APICve, error) {
	cveList := []string{}
	cveMap := make(map[string]APICve)

	client := &api.Client{HTTPClient: httpClient}
	totalPages := 9999

	// Vmaas indexes pages from 1
	for page := 1; page <= totalPages; page++ {
		vmaasRequest := APICveRequest{
			Page:             page,
			CveList:          []string{".*"},
			PageSize:         PageSize,
			RhOnly:           true,
			ErrataAssociated: true,
		}
		vmaasResponse := APICveResponse{}

		statusCode, err := client.RetryRequest(http.MethodPost, VmaasCvesURL, &vmaasRequest, &vmaasResponse)
		if err != nil {
			vmaasRequestError.WithLabelValues(VmaasCvesURL, http.MethodPost, strconv.Itoa(statusCode)).Inc()
			logger.Warningf("Request %s %s failed: statusCode=%d, err=%s", http.MethodPost, VmaasCvesURL, statusCode, err)
			return cveList, cveMap, err
		}

		for cveName, cveMD := range vmaasResponse.CveList {
			cveList = append(cveList, cveName)
			cveMap[cveName] = cveMD
		}

		totalPages = vmaasResponse.Pages
		logger.Infof("Fetched VMAAS cve list: cves=%d, page=%d/%d", len(cveMap), page, vmaasResponse.Pages)
	}

	sort.Strings(cveList)

	return cveList, cveMap, nil
}
