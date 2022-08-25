package amsclient

import (
	"fmt"

	"gorm.io/gorm"

	"app/base/models"
)

const (
	NAString = "N/A"
)

func EmptyToNA(input string) string {
	if input == "" {
		return NAString
	}
	return input
}

func DBSyncClusterDetails(conn *gorm.DB, accountID int64, clusterInfoMap map[string]ClusterInfo) error {
	// Query all clusters in DB for given account
	clusterRows := []models.Cluster{}
	if err := conn.Where("account_id = ?", accountID).Order("id").Find(&clusterRows).Error; err != nil {
		return err
	}

	for _, clusterRow := range clusterRows {
		if clusterInfo, ok := clusterInfoMap[clusterRow.UUID.String()]; ok {
			// Build provider string including region
			providerStr := EmptyToNA(clusterInfo.Provider)
			if providerStr != NAString && clusterInfo.Region != "" {
				providerStr = fmt.Sprintf("%s (%s)", providerStr, clusterInfo.Region)
			}
			if clusterRow.DisplayName != clusterInfo.DisplayName ||
				clusterRow.Status != EmptyToNA(clusterInfo.Status) ||
				clusterRow.Type != EmptyToNA(clusterInfo.Type) ||
				clusterRow.Version != EmptyToNA(clusterInfo.Version) ||
				clusterRow.Provider != providerStr {
				clusterRow.DisplayName = clusterInfo.DisplayName
				clusterRow.Status = EmptyToNA(clusterInfo.Status)
				clusterRow.Type = EmptyToNA(clusterInfo.Type)
				clusterRow.Version = EmptyToNA(clusterInfo.Version)
				clusterRow.Provider = providerStr
				if err := conn.Save(&clusterRow).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}
