package middlewares

import (
	"app/base/models"
	"app/base/utils"
	"app/test"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redhatinsights/identity"
	"github.com/stretchr/testify/assert"
)

var (
	identityMock = identity.XRHID{Identity: identity.Identity{
		Type:          "User",
		AccountNumber: "13",
		User: identity.User{
			Username:  "unit@test.com",
			Email:     "unit@test.com",
			FirstName: "Unit",
			LastName:  "Test",
			Active:    true,
			OrgAdmin:  false,
			Internal:  false,
			Locale:    "en_US",
			UserID:    "1337",
		},
		Internal: identity.Internal{},
	}}
)

var authenticator gin.HandlerFunc

func TestAuthenticatorValid(t *testing.T) {
	// valid account in db
	identityMock.Identity.OrgID = "013"
	identityMock.Identity.Internal.OrgID = "013"
	buf, _ := json.Marshal(identityMock)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test", nil)
	ctx.Request.Header.Set("x-rh-identity", base64.StdEncoding.EncodeToString(buf))
	authenticator(ctx)

	assert.Equal(t, int64(13), ctx.GetInt64("account_id"), "org_id must be translated to the acc ID")
	assert.Equal(t, identityMock.Identity.OrgID, ctx.GetString("org_id"), "org_id must be taken from identity header")
}

func TestAuthenticatorNonExisting(t *testing.T) {
	// Non existing account in db
	identityMock.Identity.OrgID = "1337"
	identityMock.Identity.Internal.OrgID = "1337"
	buf, _ := json.Marshal(identityMock)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test", nil)
	ctx.Request.Header.Set("x-rh-identity", base64.StdEncoding.EncodeToString(buf))
	authenticator(ctx)

	assert.Equal(t, int64(-1), ctx.GetInt64("account_id"), "Non-existing org_id must be translated to the account ID '-1'")
	assert.Equal(t, identityMock.Identity.OrgID, ctx.GetString("org_id"), "org_id must be taken from identity header")
}

func TestAuthenticatorInvalid(t *testing.T) {
	header := []byte("{\"wrong\":\"header\"}")

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test", nil)
	ctx.Request.Header.Set("x-rh-identity", base64.StdEncoding.EncodeToString(header))
	authenticator(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusBadRequest, "Must be 400, wrong request")
}

func TestAuthenticatorEmptyNumber(t *testing.T) {
	header := []byte("{\"org_id\":null}")

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test", nil)
	ctx.Request.Header.Set("x-rh-identity", base64.StdEncoding.EncodeToString(header))
	authenticator(ctx)

	assert.Equal(t, ctx.Writer.Status(), http.StatusBadRequest, "Must be 400, wrong request")
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
	authenticator = Authenticate(db)
	os.Exit(m.Run())
}
