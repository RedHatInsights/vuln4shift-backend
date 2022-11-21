package api

import (
	"encoding/base64"

	"github.com/pkg/errors"
)

type GithubRepoAPIResponse struct {
	Type     string `json:"type"`
	Encoding string `json:"encoding"`
	Name     string `json:"name"`
	SHA      string `json:"sha"`
	Content  string `json:"content"`
	Size     int    `json:"size"` // Max 100MB
}

func (r *GithubRepoAPIResponse) GetContents() (result []byte, err error) {
	switch r.Encoding {
	case "base64":
		result, err = base64.StdEncoding.DecodeString(r.Content)
	default:
		err = errors.Errorf("unsupported encoding: %s", r.Encoding)
	}

	return result, err
}
