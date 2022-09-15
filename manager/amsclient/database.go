package amsclient

import (
	"fmt"

	"github.com/jackc/pgtype"
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

func DBSyncClusterDetails(conn *gorm.DB, accountID int64, clusterInfoMap map[string]ClusterInfo) ([]string, map[string]struct{}, map[string]struct{}, map[string]struct{}, error) {
	clusterIDs := []string{}
	clusterStatuses := map[string]struct{}{}
	clusterVersions := map[string]struct{}{}
	clusterProviders := map[string]struct{}{}

	// Query all clusters in DB for given account
	clusterRows := []models.Cluster{}
	if err := conn.Where("account_id = ?", accountID).Order("id").Find(&clusterRows).Error; err != nil {
		return nil, nil, nil, nil, err
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
				// workaround to not encode undefined value
				if clusterRow.Workload.Status == pgtype.Undefined {
					clusterRow.Workload.Status = pgtype.Null
				}
				if err := conn.Save(&clusterRow).Error; err != nil {
					return nil, nil, nil, nil, err
				}
			}
			clusterIDs = append(clusterIDs, clusterInfo.ID)
			clusterStatuses[EmptyToNA(clusterInfo.Status)] = struct{}{}
			clusterVersions[EmptyToNA(clusterInfo.Version)] = struct{}{}
			clusterProviders[providerStr] = struct{}{}
		}
	}

	return clusterIDs, clusterStatuses, clusterVersions, clusterProviders, nil
}
