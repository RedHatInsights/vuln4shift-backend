package test

import (
	"app/base/models"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAllClusters(t *testing.T) (clusters []models.ClusterLight) {
	result := DB.Model(models.ClusterLight{}).Order("id").Scan(&clusters)
	assert.Nil(t, result.Error)
	assert.True(t, len(clusters) > 0)
	return clusters
}

func GetAccountClusters(t *testing.T, id int64) (clusters []models.ClusterLight) {
	result := DB.Model(models.ClusterLight{}).
		Order("id").
		Where("account_id = ?", id).
		Scan(&clusters)
	assert.Nil(t, result.Error)
	assert.True(t, len(clusters) > 0)
	return clusters
}

func GetCluster(t *testing.T, id int64) (cluster models.ClusterLight) {
	result := DB.Model(models.ClusterLight{}).Where("id = ?", id).Scan(&cluster)
	assert.Nil(t, result.Error)
	return cluster
}

func CheckClusterMetaSlices(t *testing.T, ep, ap, es, as, ev, av []string) {
	sort.Strings(ep)
	sort.Strings(es)
	sort.Strings(ev)
	sort.Strings(ap)
	sort.Strings(as)
	sort.Strings(av)

	assert.Equal(t, ep, ap)
	assert.Equal(t, es, as)
	assert.Equal(t, ev, av)
}

func CheckClustersMeta(t *testing.T, meta interface{}, providers, statuses, versions map[string]bool) {
	// Actual meta
	am := meta.(map[string]interface{})
	// Expected meta slices without duplicates
	ep := GetMetaKeys(providers)
	es := GetMetaKeys(statuses)
	ev := GetMetaKeys(versions)
	// Actual meta slices
	ap := GetMetaStringSlice(am["cluster_providers_all"])
	as := GetMetaStringSlice(am["cluster_statuses_all"])
	av := GetMetaStringSlice(am["cluster_versions_all"])

	CheckClusterMetaSlices(t, ep, ap, es, as, ev, av)
}

func CheckClusterDetails(t *testing.T, ep, es, ev map[string]bool, ap, as, av map[string]struct{}) {
	CheckClusterMetaSlices(t,
		GetMetaKeys(ep), GetClusterDetailKeys(ap),
		GetMetaKeys(es), GetClusterDetailKeys(as),
		GetMetaKeys(ev), GetClusterDetailKeys(av))
}
