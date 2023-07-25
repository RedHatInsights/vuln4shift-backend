package vmsync

import (
	"app/base/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	DB       *gorm.DB
	dbCveMap map[string]models.Cve
)

func prepareDbCvesMap() error {
	var cveRows []models.Cve
	if err := DB.Order("name").Find(&cveRows).Error; err != nil {
		return err
	}
	dbCveMap = make(map[string]models.Cve, len(cveRows))
	for _, cve := range cveRows {
		dbCveMap[cve.Name] = cve
	}
	return nil
}

func insertUpdateCves(toSyncCves []models.Cve, tx *gorm.DB) error {
	logger.Debugf("CVEs to insert/update: %d", len(toSyncCves))

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).CreateInBatches(toSyncCves, BatchSize).Error
}

// deleteNotAffectingCves deletes CVEs not affecting any cluster.
func deleteNotAffectingCves(tx *gorm.DB, cveIds []int64) (int64, error) {
	var toDeleteIds []int64
	err := tx.Table("cve").
		Select("cve.id").
		Joins("full outer join image_cve as ic on cve.id = ic.cve_id").
		Where("ic.image_id is NULL and cve.id in (?)", cveIds).
		Scan(&toDeleteIds).Error
	if err != nil {
		return 0, err
	}
	logger.Infof("CVEs to delete: %d", len(toDeleteIds))

	if err := tx.Where("cve_id in ?", toDeleteIds).Delete(&models.AccountCveCache{}).Error; err != nil {
		return 0, err
	}
	if err := tx.Where("cve_id in ?", toDeleteIds).Delete(&models.ClusterCveCache{}).Error; err != nil {
		return 0, err
	}
	if err := tx.Where("cve_id in ?", toDeleteIds).Delete(&models.ImageCve{}).Error; err != nil {
		return 0, err
	}
	if err := tx.Where("id in ?", toDeleteIds).Delete(&models.Cve{}).Error; err != nil {
		return 0, err
	}

	return int64(len(toDeleteIds)), nil
}
