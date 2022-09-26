package clusters

import (
	"app/base/logging"
	"app/base/models"
	"app/base/utils"
	"app/manager/middlewares"
	"app/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

var (
	testFilterer   = middlewares.Filterer()
	testController Controller
)

func callGetClusters(t *testing.T, accountID int64, expectedStatus int) *httptest.ResponseRecorder {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	ctx, w := test.MockGinRequest(header, "GET", nil, gin.Params{}, url.Values{})
	ctx.Set("account_id", accountID)

	testFilterer(ctx)
	testController.GetClusters(ctx)

	assert.Equal(t, expectedStatus, w.Code)
	return w
}

func TestGetClusters(t *testing.T) {
	allAccounts := test.GetAccounts(t)
	for _, account := range allAccounts {
		var resp GetClustersResponse
		w := callGetClusters(t, account.ID, http.StatusOK)
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &resp))

		expectedClusters := test.GetAccountClusters(t, account.ID)
		assert.Equal(t, len(expectedClusters), len(resp.Data))
		for i, ec := range expectedClusters {
			// Actual cluster.
			ac := resp.Data[i]

			assert.Equal(t, ec.UUID.String(), ac.UUID)
			assert.Equal(t, ec.Status, ac.Status)
			if ec.Type != "" {
				assert.Equal(t, ec.Type, ac.Type)
			}
			assert.Equal(t, ec.UUID.String(), ac.DisplayName)
			assert.Equal(t, ec.Version, ac.Version)
			assert.Equal(t, ec.LastSeen.UTC(), ac.LastSeen.UTC())

			// Expected severities count.
			es := test.GetCvesTypeCount(test.GetClusterCves(t, ec.ID))
			assert.Equal(t, es[models.Critical], *ac.Severities.CriticalCount)
			assert.Equal(t, es[models.Important], *ac.Severities.ImportantCount)
			assert.Equal(t, es[models.Moderate], *ac.Severities.ModerateCount)
			assert.Equal(t, es[models.Low], *ac.Severities.LowCount)
		}
	}
}

func TestMain(m *testing.M) {
	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		panic(err)
	}

	test.DB = db
	err = test.ResetDB()
	if err != nil {
		panic(err)
	}

	logger, err := logging.CreateLogger("DEBUG")
	if err != nil {
		panic(err)
	}

	testController.Conn = test.DB
	testController.Logger = logger

	os.Exit(m.Run())
}
