package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"app/base/models"
	"app/base/utils"
	"app/manager/base"

	"github.com/gin-gonic/gin"
	"github.com/redhatinsights/identity"
	"gorm.io/gorm"
)

func getOrgID(id *identity.XRHID) (string, error) {
	// FIXME: keycloak in ephemeral currently sets only the internal org_id,
	// this may be fixed in future, so we'll not need to check both
	if id.Identity.OrgID != "" {
		return id.Identity.OrgID, nil
	} else if id.Identity.Internal.OrgID != "" {
		return id.Identity.Internal.OrgID, nil
	}
	return "", errors.New("x-rh-identity header has an invalid or missing org_id")
}

// Authenticate checks if user
// has x-rh-identity sent and upserts
// the user to the db
func Authenticate(db *gorm.DB) gin.HandlerFunc {
	logger, err := utils.CreateLogger(utils.Cfg.LoggingLevel)
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

		orgID, err := getOrgID(&id)
		if err != nil {
			logger.Debug(fmt.Sprintf("Invalid x-rh-identity obtained: %s, %s", idS, err.Error()))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, base.BuildErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		var acc models.Account
		res := db.Where("org_id = ?", orgID).Find(&acc)

		if res.RowsAffected > 0 {
			ctx.Set("account_id", acc.ID)
			ctx.Set("org_id", orgID)
		} else {
			// set non-existing account_id, so account with empty systems, can still get response with empty data
			ctx.Set("account_id", int64(-1))
			ctx.Set("org_id", orgID)
		}
	}
}
