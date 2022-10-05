package test

import (
	"app/base/ams"
)

type AMSClientMock struct {
	ClustersResponse map[string]ams.ClusterInfo
	ClusterResponse  ams.ClusterInfo
}

func (c *AMSClientMock) GetClustersForOrganization(orgID string) (
	map[string]ams.ClusterInfo,
	error,
) {
	return c.ClustersResponse, nil
}

func (c *AMSClientMock) GetSingleClusterInfoForOrganization(orgID string, clusterID string) (
	ams.ClusterInfo, error,
) {
	return c.ClusterResponse, nil
}
