package cleaner_test

import (
	"app/base/models"
	"app/base/utils"
	"app/cleaner"
	"app/test"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var clusterCleaner *cleaner.ClusterCleaner

func TestRemoveEmpty(t *testing.T) {
	var ogClusterCnt int64
	var afterClusterCnt int64

	test.DB.Model(&models.Cluster{}).Count(&ogClusterCnt)

	err := clusterCleaner.Clean()
	test.DB.Model(&models.Cluster{}).Count(&afterClusterCnt)

	assert.Equal(t, nil, err, "Error from cleaner should be nil")
	assert.Equal(t, ogClusterCnt, afterClusterCnt, "No expired clusters in DB, expecting same count of clusters")
}

func TestRemove(t *testing.T) {
	var ogClusterCnt int64
	var afterClusterCnt int64

	expiredTimestamp := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	expiredCluster := models.Cluster{UUID: uuid.New(), AccountID: 13, LastSeen: expiredTimestamp, Status: "s", Version: "v"}
	test.DB.Save(&expiredCluster)
	test.DB.Model(&models.Cluster{}).Count(&ogClusterCnt)

	err := clusterCleaner.Clean()
	test.DB.Model(&models.Cluster{}).Count(&afterClusterCnt)

	assert.Equal(t, nil, err, "Error from cleaner should be nil")
	assert.Equal(t, ogClusterCnt-1, afterClusterCnt, "One cluster needs to be purged")
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
	// 30 retention days
	clusterCleaner, err = cleaner.NewClusterCleaner(30)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
