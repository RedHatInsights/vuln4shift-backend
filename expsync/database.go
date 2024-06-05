package expsync

import (
	"app/base/models"
	"encoding/json"

	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func updateExploitsMetadata(tx *gorm.DB, exploits map[CVE][]ExploitMetadata) (int64, error) {
	if len(exploits) == 0 {
		return 0, nil
	}

	var totalRowsAffected int64
	for cve, metadata := range exploits {
		exploitBytes, err := json.Marshal(metadata)
		if err != nil {
			return 0, err
		}
		res := tx.Exec("UPDATE cve SET exploit_data = ? WHERE name = ?", exploitBytes, cve)
		if e := res.Error; e != nil {
			return 0, e
		}
		totalRowsAffected = totalRowsAffected + res.RowsAffected
	}

	return totalRowsAffected, nil
}

func getCvesWithExploitMetadata(tx *gorm.DB) (cves []models.Cve, err error) {
	return cves, tx.Model(models.Cve{}).Where("exploit_data IS NOT NULL").Find(&cves).Error
}

func removeExploitData(tx *gorm.DB, cves []models.Cve) (int64, error) {
	if len(cves) == 0 {
		return 0, nil
	}

	res := tx.Model(cves).UpdateColumn("exploit_data", nil)
	return res.RowsAffected, res.Error
}
