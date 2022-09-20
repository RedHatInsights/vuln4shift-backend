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

func DeleteCvesByID(t *testing.T, ids ...int64) {
	result := DB.Delete(&models.Cve{}, ids)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(len(ids)), result.RowsAffected)
}
