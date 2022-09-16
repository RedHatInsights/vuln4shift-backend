package test

import (
	"app/base/models"
	"testing"

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
