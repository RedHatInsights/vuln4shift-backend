package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"app/base/logging"
	"app/base/utils"
)

var (
	logger  *logrus.Logger
	retries = utils.Cfg.APIRetries
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient HTTPClient
}

func init() {
	var err error
	logger, err = logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func (c *Client) Request(method, url string, requestPtr, responsePtr interface{}) (int, error) {
	body := &bytes.Buffer{}
	if requestPtr != nil {
		err := json.NewEncoder(body).Encode(requestPtr)
		if err != nil {
			return 0, err
		}
	}

	request, err := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return 0, err
	}

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return 0, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, err
	}

	err = json.Unmarshal(bodyBytes, responsePtr)
	if err != nil {
		return response.StatusCode, err
	}

	return response.StatusCode, nil
}

func (c *Client) RetryRequest(method string, url string, requestPtr interface{}, responsePtr interface{}) (int, error) {
	var statusCode int
	var err error
	for i := 0; i < retries; i++ {
		statusCode, err = c.Request(method, url, requestPtr, responsePtr)
		if statusCode >= 200 && statusCode <= 299 && err == nil {
			return statusCode, err
		}
		if i < (retries - 1) {
			logger.Debugf("Request %s %s failed, retrying: statusCode=%d, err=%s", method, url, statusCode, err)
		}
	}
	return statusCode, err
}
