package test

import (
	"app/base/models"
	"testing"

	"gorm.io/gorm/clause"

	"github.com/stretchr/testify/assert"
)

func GetAllCves(t *testing.T) []models.Cve {
	var cves []models.Cve
	assert.Nil(t, DB.Model(models.Cve{}).Scan(&cves).Error)
	return cves
}

func GetNonAffectingCves(t *testing.T) []models.Cve {
	var cves []models.Cve
	assert.Nil(t, DB.Model(models.Cve{}).
		Joins("full outer join image_cve as ic on id = ic.cve_id").
		Where("ic.image_id is NULL").
		Scan(&cves).Error)
	return cves
}

func InsertCve(t *testing.T, cve models.Cve) int64 {
	result := DB.Model(models.Cve{}).Create(&cve)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(1), result.RowsAffected)
	return cve.ID
}

func UpsertCves(t *testing.T, cves []models.Cve) {
	result := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: false,
		UpdateAll: true,
	}).Create(&cves)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(len(cves)), result.RowsAffected)
}

func GetCveByID(t *testing.T, id int64) *models.Cve {
	var cve *models.Cve
	assert.Nil(t, DB.Model(models.Cve{}).Where("id", id).Scan(&cve).Error)
	return cve
}

func GetCvesByName(t *testing.T, names ...string) (cves []models.Cve) {
	assert.Nil(t, DB.Model(models.Cve{}).Where("name in (?)", names).Order("id").Scan(&cves).Error)
	return cves
}

func DeleteCvesByID(t *testing.T, ids ...int64) {
	result := DB.Delete(&models.Cve{}, ids)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(len(ids)), result.RowsAffected)
}

func GetClusterCves(t *testing.T, id int64) (cves []models.Cve) {
	assert.Nil(t, DB.Model(models.Cve{}).
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Group("cve.id").Order("cve.id").
		Where("cluster.id = ?", id).
		Scan(&cves).Error)
	return cves
}

func GetCvesTypeCount(cves []models.Cve) map[models.Severity]int64 {
	res := make(map[models.Severity]int64)
	for _, cve := range cves {
		res[cve.Severity]++
	}
	return res
}

func GetAccountCves(t *testing.T, id int64) (cves []models.Cve) {
	assert.Nil(t, DB.Model(models.Cve{}).
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Group("cve.id").Order("cve.id").
		Where("cluster.account_id = ?", id).
		Scan(&cves).Error)
	return cves
}

func GetAccountCvesForClusters(t *testing.T, id int64, clusterIDs []int64) (cves []models.Cve) {
	assert.Nil(t, DB.Model(models.Cve{}).
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Group("cve.id").Order("cve.id").
		Where("cluster.account_id = ? AND cluster.id = (?)", id, clusterIDs).
		Scan(&cves).Error)
	return cves
}

func GetImagesExposed(t *testing.T, accountID, cveID int64) int64 {
	var imagesExposed int64
	assert.Nil(t, DB.Model(models.Cve{}).
		Select("COUNT(repository_image.repository_id)").
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Joins("JOIN repository_image ON cluster_image.image_id = repository_image.image_id").
		Group("cve.id").
		Where("cluster.account_id = ? AND cve.id = ?", accountID, cveID).
		Scan(&imagesExposed).Error)
	return imagesExposed
}

func GetImagesExposedLimitClusters(t *testing.T, accountID, cveID int64, clusterIDs []int64) int64 {
	var imagesExposed int64
	assert.Nil(t, DB.Model(models.Cve{}).
		Select("COUNT(repository_image.repository_id)").
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Joins("JOIN repository_image ON cluster_image.image_id = repository_image.image_id").
		Group("cve.id").
		Where("cluster.account_id = ? AND cve.id = ? AND cluster.id = (?)", accountID, cveID, clusterIDs).
		Scan(&imagesExposed).Error)
	return imagesExposed
}

func GetClustersExposed(t *testing.T, accountID, cveID int64) int64 {
	var clustersExposed int64
	assert.Nil(t, DB.Model(models.Cve{}).
		Select("COUNT(DISTINCT cluster_image.cluster_id)").
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Group("cve.id").
		Where("cluster.account_id = ? AND cve.id = ?", accountID, cveID).
		Scan(&clustersExposed).Error)
	return clustersExposed
}

func GetClustersExposedLimitClusters(t *testing.T, accountID, cveID int64, clusterIDs []int64) int64 {
	var clustersExposed int64
	assert.Nil(t, DB.Model(models.Cve{}).
		Select("COUNT(DISTINCT cluster_image.cluster_id)").
		Joins("JOIN image_cve ON cve.id = image_cve.cve_id").
		Joins("JOIN cluster_image ON image_cve.image_id = cluster_image.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Group("cve.id").
		Where("cluster.account_id = ? AND cve.id = ? AND cluster.id = (?)", accountID, cveID, clusterIDs).
		Scan(&clustersExposed).Error)
	return clustersExposed
}

func GetExposedClusters(t *testing.T, accountID, cveID int64) (clusters []models.ClusterLight) {
	assert.Nil(t, DB.Model(models.ClusterLight{}).
		Joins("JOIN cluster_image on cluster.id = cluster_image.cluster_id").
		Joins("JOIN image_cve on cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cve on image_cve.cve_id = cve.id").
		Group("cluster.id").
		Where("cluster.account_id = ? AND cve.id = ?", accountID, cveID).
		Scan(&clusters).Error)
	return clusters
}
