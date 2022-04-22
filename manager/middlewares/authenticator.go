package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"app/base/logging"
	"app/base/models"
	"app/base/utils"
	"app/manager/base"

	"github.com/gin-gonic/gin"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"gorm.io/gorm"
)

func checkIdentity(id *identity.XRHID) error {
	if id.Identity.Type == "Associate" && id.Identity.AccountNumber == "" {
		return nil
	}
	if id.Identity.AccountNumber == "" || id.Identity.AccountNumber == "-1" {
		return errors.New("x-rh-identity header has an invalid or missing account number")
	}
	if id.Identity.OrgID == "" || id.Identity.Internal.OrgID == "" {
		return errors.New("x-rh-identity header has an invalid or missing org_id")
	}
	if id.Identity.Type == "" {
		return errors.New("x-rh-identity header is missing type")
	}
	return nil
}

// Authenticate checks if user
// has x-rh-identity sent and upserts
// the user to the db
func Authenticate(db *gorm.DB) gin.HandlerFunc {
	logger, err := logging.CreateLogger(utils.GetEnv("LOGGING_LEVEL", "DEBUG"))
	if err != nil {
		panic("Invalid LOGGING_LEVEL enviroment variable set")
	}

	return func(ctx *gin.Context) {
		idS := ctx.GetHeader("x-rh-identity")

		idRaw, err := base64.StdEncoding.DecodeString(idS)
		if err != nil {
			logger.Debug(fmt.Sprintf("Invalid x-rh-identity obtained: %s, %s", idS, err.Error()))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, "invalid B64 x-rh-identity"))
			return
		}

		var id identity.XRHID
		err = json.Unmarshal(idRaw, &id)
		if err != nil {
			logger.Debug(fmt.Sprintf("Invalid x-rh-identity obtained: %s, %s", idS, err.Error()))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, "cannot parse x-rh-identity"))
			return
		}

		err = checkIdentity(&id)
		if err != nil {
			logger.Debug(fmt.Sprintf("Invalid x-rh-identity obtained: %s, %s", idS, err.Error()))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		var acc models.Account
		res := db.Where("org_id = ?", id.Identity.OrgID).Find(&acc)

		if res.RowsAffected > 0 {
			ctx.Set("account_id", acc.ID)
			ctx.Set("org_id", id.Identity.OrgID)
		} else {
			// set non-existing account_id, so account with empty systems, can still get response with empty data
			ctx.Set("account_id", int64(-1))
			ctx.Set("org_id", id.Identity.OrgID)
		}
	}
}
