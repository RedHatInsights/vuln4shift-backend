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

func DBFetchClusterDetails(conn *gorm.DB, ams AMSClient, accountID int64, orgID string, sync bool, cveName *string) ([]string, map[string]struct{}, map[string]struct{}, map[string]struct{}, error) {
	clusterIDs := []string{}
	clusterStatuses := map[string]struct{}{}
	clusterVersions := map[string]struct{}{}
	clusterProviders := map[string]struct{}{}

	// Query all clusters in DB for given account
	clusterRows := []models.ClusterLight{}
	query := conn.Where("cluster.account_id = ?", accountID).Order("cluster.id")
	if cveName != nil {
		query = query.
			Joins("JOIN cluster_image ON cluster.id = cluster_image.cluster_id").
			Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
			Joins("JOIN cve ON image_cve.cve_id = cve.id").
			Where("cve.name = ?", *cveName).
			Distinct()
	}
	if err := query.Find(&clusterRows).Error; err != nil {
		return nil, nil, nil, nil, err
	}

	if sync {
		clusterInfoMap, err := ams.GetClustersForOrganization(orgID)
		if err != nil {
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
	} else {
		// Return current values from DB is sync is set to false
		for _, clusterRow := range clusterRows {
			clusterIDs = append(clusterIDs, clusterRow.UUID.String())
			clusterStatuses[clusterRow.Status] = struct{}{}
			clusterVersions[clusterRow.Version] = struct{}{}
			clusterProviders[clusterRow.Provider] = struct{}{}
		}
	}

	return clusterIDs, clusterStatuses, clusterVersions, clusterProviders, nil
}
