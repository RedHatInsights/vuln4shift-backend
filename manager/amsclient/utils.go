package amsclient

import (
	"fmt"
	"strings"
)

// generateSearchParameter generates a search string for given org_id and desired statuses
func generateSearchParameter(orgID string, allowedStatuses, disallowedStatuses []string, clusterSearch string) string {
	searchQuery := fmt.Sprintf("organization_id is '%s' and cluster_id != ''", orgID)

	if len(allowedStatuses) > 0 {
		clusterIDQuery := " and status in ('" + strings.Join(allowedStatuses, "','") + "')"
		searchQuery += clusterIDQuery
	}

	if len(disallowedStatuses) > 0 {
		clusterIDQuery := " and status not in ('" + strings.Join(disallowedStatuses, "','") + "')"
		searchQuery += clusterIDQuery
	}

	if len(clusterSearch) > 0 {
		clusterIDQuery := " and (display_name ilike '%" + clusterSearch + "%' or external_cluster_id ilike '%" + clusterSearch + "%')"
		searchQuery += clusterIDQuery
	}

	return searchQuery
}
