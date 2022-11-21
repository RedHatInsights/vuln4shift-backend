package expsync

import (
	"app/base/api"
	"app/base/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

var (
	exploitFileURL                = utils.Cfg.ExploitFileURL
	httpClient     api.HTTPClient = &http.Client{}
)

type ExploitMetadata struct {
	Date      string `json:"date"`
	Source    string `json:"source"`
	Reference string `json:"reference"`
}

type CVE string

type ExploitAPIResponse map[CVE][]ExploitMetadata

// getAPIExploits gets exploit metadata file from VMaaS-assets resource.
func getAPIExploits() (map[CVE][]ExploitMetadata, error) {
	header := http.Header{}
	header.Set("Accept", "application/vnd.github+json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", utils.Cfg.GitToken))

	client := &api.Client{HTTPClient: httpClient, Header: &header}

	const method = http.MethodGet

	var githubAPIResp api.GithubRepoAPIResponse
	code, err := client.RetryRequest(method, exploitFileURL, nil, &githubAPIResp)
	if err != nil {
		exploitsRequestError.WithLabelValues(exploitFileURL, method, strconv.Itoa(code)).Inc()
		logger.Warningf("Request %s %s failed: statusCode=%d, err=%s", method, exploitFileURL, code, err)
		return nil, err
	}

	if code == http.StatusUnauthorized {
		return nil, errors.New("invalid API token provided")
	}
	if !api.IsSuccessCode(code) {
		return nil, errors.Errorf("unexpected status code received: %d", code)
	}

	contents, err := githubAPIResp.GetContents()
	if err != nil {
		logger.Errorf("Decoding github API response content failed: %s", err.Error())
		return nil, errors.Wrap(err, "failed to get content from github API response")
	}

	var exploitsResp ExploitAPIResponse
	err = json.Unmarshal(contents, &exploitsResp)
	if err != nil {
		logger.Errorf("Unexpected data received from github API: %s", string(contents))
		return nil, errors.Wrap(err, "unmarshalling github API response to exploits metadata failed")
	}

	return exploitsResp, nil
}
