package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	githubAPIResp = `{
		"type": "file",
		"encoding": "base64",
		"size": 5362,
		"name": "exploits.json",
		"path": "exploits.json",
		"content": "ewoJCSJDVkUtMjAyMi0wMDA0IjogWwoJCQl7CgkJCQkiZGF0ZSI6ICIyMDIyLTA4LTE4IiwKCQkJCSJyZWZlcmVuY2UiOiAiTi9BIiwKCQkJCSJzb3VyY2UiOiAiQ0lTQSIKCQkJfQoJCV0sCgkJIkNWRS0yMDIyLTAwMDUiOiBbCgkJCXsKCQkJCSJkYXRlIjogIjIwMjItMTAtMjgiLAoJCQkJInJlZmVyZW5jZSI6ICJOL0EiLAoJCQkJInNvdXJjZSI6ICJDSVNBIgoJCQl9CgkJXQoJfQ==\n",
		"sha": "3d21ec53a331a6f037a91c368710b99387d012c1"
	}`

	expectedContents = `{
		"CVE-2022-0004": [
			{
				"date": "2022-08-18",
				"reference": "N/A",
				"source": "CISA"
			}
		],
		"CVE-2022-0005": [
			{
				"date": "2022-10-28",
				"reference": "N/A",
				"source": "CISA"
			}
		]
	}`

	githubAPIRespUnknownEncoding = `{"encoding": "unknown"}`
)

func TestGetContents(t *testing.T) {
	var resp GithubRepoAPIResponse
	assert.Nil(t, json.Unmarshal([]byte(githubAPIResp), &resp))

	content, err := resp.GetContents()
	assert.Nil(t, err)

	assert.Equal(t, expectedContents, string(content))
	assert.Equal(t, "file", resp.Type)
	assert.Equal(t, "exploits.json", resp.Name)
}

func TestGetContentsUnknownEncoding(t *testing.T) {
	var resp GithubRepoAPIResponse
	assert.Nil(t, json.Unmarshal([]byte(githubAPIRespUnknownEncoding), &resp))

	content, err := resp.GetContents()
	assert.Nil(t, content)
	assert.Equal(t, "unsupported encoding: unknown", err.Error())
}
