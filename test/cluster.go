package test

import (
	"app/base/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAllClusters(t *testing.T) (clusters []models.Cluster) {
	result := DB.Model(models.Cluster{}).Order("id").Scan(&clusters)
	assert.Nil(t, result.Error)
	assert.True(t, len(clusters) > 0)
	return clusters
}

func GetAccountClusters(t *testing.T, id int64) (clusters []models.Cluster) {
	result := DB.Model(models.Cluster{}).
		Order("id").
		Where("account_id = ?", id).
		Scan(&clusters)
	assert.Nil(t, result.Error)
	assert.True(t, len(clusters) > 0)
	return clusters
}

func GetCluster(t *testing.T, id int64) (cluster models.Cluster) {
	result := DB.Model(models.Cluster{}).Where("id = ?", id).Scan(&cluster)
	assert.Nil(t, result.Error)
	return cluster
}
