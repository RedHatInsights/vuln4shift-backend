package expsync

import (
	"app/base/models"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func updateExploitsMetadata(tx *gorm.DB, exploits map[CVE][]ExploitMetadata) (int64, error) {
	if len(exploits) == 0 {
		return 0, nil
	}

	var argList string
	for cve, metadata := range exploits {
		exploitBytes, err := json.Marshal(metadata)
		if err != nil {
			return 0, err
		}
		argList = fmt.Sprintf("%s ('%s', '%s'::jsonb),", argList, string(cve), exploitBytes)
	}

	// Trim preceding whitespace and trailing comma.
	argList = argList[1 : len(argList)-1]

	statement := fmt.Sprintf("UPDATE cve SET exploit_data = exploits.data FROM (VALUES %s) AS exploits(cve, data) WHERE cve.name = exploits.cve", argList)
	res := tx.Exec(statement)
	if e := res.Error; e != nil {
		return 0, e
	}

	return res.RowsAffected, nil
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
