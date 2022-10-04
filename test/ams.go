package test

import (
	"app/manager/amsclient"
)

type AMSClientMock struct {
	ClustersResponse map[string]amsclient.ClusterInfo
	ClusterResponse  amsclient.ClusterInfo
}

func (c *AMSClientMock) GetClustersForOrganization(orgID string) (
	map[string]amsclient.ClusterInfo,
	error,
) {
	return c.ClustersResponse, nil
}

func (c *AMSClientMock) GetSingleClusterInfoForOrganization(orgID string, clusterID string) (
	amsclient.ClusterInfo, error,
) {
	return c.ClusterResponse, nil
}
