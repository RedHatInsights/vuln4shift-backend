package vmsync

import (
	"app/base/models"
	"app/base/utils"

	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	dbCveMap map[string]models.Cve
)

func dbConfigure() error {
	dsn := utils.GetDbURL()
	var err error
	DB, err = models.GetGormConnection(dsn)

	if err != nil {
		return err
	}
	return nil
}

func prepareDbCvesMap() error {
	var cveRows []models.Cve
	if err := DB.Find(&cveRows).Error; err != nil {
		return err
	}
	dbCveMap = make(map[string]models.Cve, len(cveRows))
	for _, cve := range cveRows {
		dbCveMap[cve.Name] = cve
	}
	return nil
}
