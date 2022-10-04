package cves

import (
	"app/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func callGetCveDetails(t *testing.T, accountID int64, cveName string, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	params := gin.Params{{Key: "cve_name", Value: cveName}}

	ctx, w := test.MockGinRequest(header, "GET", nil, params, url.Values{})
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetCveDetails(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetCveDetails(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		cves := test.GetAccountCves(t, account.ID)
		for _, cve := range cves {
			var resp GetCveDetailsResponse
			w := callGetCveDetails(t, account.ID, cve.Name, http.StatusOK)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))
			// Actual CVE
			ac := resp.Data
			assert.Equal(t, cve.Severity, ac.Severity)
			assert.Equal(t, cve.Description, ac.Description)
			assert.Equal(t, cve.Name, ac.Name)
			assert.Equal(t, cve.Cvss2Score, test.GetFloat32PtrValue(ac.Cvss2Score))
			assert.Equal(t, cve.Cvss2Metrics, test.GetStringPtrValue(ac.Cvss2Metrics))
			assert.Equal(t, cve.Cvss3Score, test.GetFloat32PtrValue(ac.Cvss3Score))
			assert.Equal(t, cve.Cvss3Metrics, test.GetStringPtrValue(ac.Cvss3Metrics))
			assert.Equal(t, cve.RedhatURL, test.GetStringPtrValue(ac.RedhatURL))
			assert.Equal(t, test.GetUTC(cve.PublicDate), test.GetUTC(ac.PublicDate))
		}
	}
}
