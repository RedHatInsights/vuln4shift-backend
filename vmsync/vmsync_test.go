package vmsync

import (
	"app/base/models"
	"app/base/utils"
	"app/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPruneCves(t *testing.T) {
	beforeSyncCnt := len(test.GetAllCves(t))
	nonAffectingCnt := len(test.GetNonAffectingCves(t))
	assert.Nil(t, prepareDbCvesMap())
	assert.Nil(t, pruneCves())
	afterSyncCnt := len(test.GetAllCves(t))
	assert.True(t, beforeSyncCnt-afterSyncCnt == nonAffectingCnt)
}

func TestMain(m *testing.M) {
	db, err := models.GetGormConnection(utils.GetDbURL(false))
	if err != nil {
		panic(err)
	}
	test.DB = db
	DB = test.DB
	err = test.ResetDB()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
