// Logic of this package borrowed from https://github.com/RedHatInsights/insights-results-smart-proxy
package amsclient

import (
	"fmt"
	"net/http"
	"time"

	sdk "github.com/openshift-online/ocm-sdk-go"
	accMgmt "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"

	"github.com/sirupsen/logrus"

	"app/base/logging"
	"app/base/utils"
)

const (
	// strings for logging and errors
	orgNoInternalID     = "Organization doesn't have proper internal ID"
	orgMoreInternalOrgs = "More than one internal organization for the given orgID"

	// StatusDeprovisioned indicates the corresponding cluster subscription status
	StatusDeprovisioned = "Deprovisioned"
	// StatusArchived indicates the corresponding cluster subscription status
	StatusArchived = "Archived"
	// StatusReserved means the cluster has reserved resources, but isn't initialized yet.
	StatusReserved = "Reserved"
)

var (
	// DefaultStatusNegativeFilters are filters that are applied to the AMS API subscriptions query when the filters are empty
	// We are either not interested in clusters in these states (Archived, Deprovisioned) or the cluster's
	// initialization hasn't finished yet (Reserved), meaning the cluster is not ready to start sending Insights archives,
	// as it might not even have a Cluster UUID assigned yet. When the initialization succeeds or fails, the cluster's
	// state becomes either Active or Deprovisioned.
	DefaultStatusNegativeFilters = []string{StatusArchived, StatusDeprovisioned, StatusReserved}
)

// AMSClient allow us to interact the AMS API
type AMSClient interface {
	GetClustersForOrganization(string) (
		clusterInfoMap map[string]ClusterInfo,
		err error,
	)
	GetSingleClusterInfoForOrganization(string, string) (
		ClusterInfo, error,
	)
}

// amsClientImpl is an implementation of the AMSClient interface
type amsClientImpl struct {
	connection *sdk.Connection
	pageSize   int
	logger     *logrus.Logger
}

// NewAMSClient create an AMSClient from the configuration
func NewAMSClient() (AMSClient, error) {
	return NewAMSClientWithTransport(nil)
}

// NewAMSClientWithTransport creates an AMSClient from the configuration, enabling to use a transport wrapper
func NewAMSClientWithTransport(transport http.RoundTripper) (AMSClient, error) {
	logger, err := logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		panic("invalid LOGGING_LEVEL enviroment variable set")
	}
	logger.Debugln("creating amsclient...")
	builder := sdk.NewConnectionBuilder().URL(utils.Cfg.AmsAPIURL)

	if transport != nil {
		builder.TransportWrapper(func(http.RoundTripper) http.RoundTripper { return transport })
	}

	if utils.Cfg.AmsClientID != "" && utils.Cfg.AmsClientSecret != "" {
		builder = builder.Client(utils.Cfg.AmsClientID, utils.Cfg.AmsClientSecret)
	} else {
		err := fmt.Errorf("no credentials provided, cannot create the API client")
		logger.Errorln(err.Error())
		return nil, err
	}

	conn, err := builder.Build()

	if err != nil {
		logger.Errorf("unable to build the connection to AMS API: %s", err.Error())
		return nil, err
	}

	return &amsClientImpl{
		connection: conn,
		pageSize:   utils.Cfg.AmsAPIPagesize,
		logger:     logger,
	}, nil
}

// GetClustersForOrganization retrieves the clusters for a given organization using the default client
// it allows to filter the clusters by their status (statusNegativeFilter will exclude the clusters with status in that list)
// If nil is passed for filters, default filters will be applied. To select empty filters, pass an empty slice.
func (c *amsClientImpl) GetClustersForOrganization(orgID string) (
	clusterInfoMap map[string]ClusterInfo,
	err error,
) {
	c.logger.Debugf("looking up active clusters for the organization: %s", orgID)
	c.logger.Debugf("GetClustersForOrganization start, AMS client page size: %v", c.pageSize)

	tStart := time.Now()

	internalOrgID, err := c.GetInternalOrgIDFromExternal(orgID)
	if err != nil {
		return
	}

	searchQuery := generateSearchParameter(internalOrgID, DefaultStatusNegativeFilters)
	subscriptionListRequest := c.connection.AccountsMgmt().V1().Subscriptions().List()

	clusterInfoMap, err = c.executeSubscriptionListRequest(subscriptionListRequest, searchQuery)
	if err != nil {
		c.logger.Errorf("GetClustersForOrganization err=%s, org_id=%s", err.Error(), orgID)
		return
	}

	c.logger.Debugf("GetClustersForOrganization from AMS API took %s, org_id=%s", time.Since(tStart), orgID)
	return
}

func (c *amsClientImpl) GetSingleClusterInfoForOrganization(orgID string, clusterID string) (
	clusterInfo ClusterInfo, err error,
) {
	tStart := time.Now()

	internalOrgID, err := c.GetInternalOrgIDFromExternal(orgID)
	if err != nil {
		return
	}

	searchQuery := fmt.Sprintf("organization_id = '%s' and external_cluster_id = '%s'", internalOrgID, clusterID)

	subscriptionListRequest := c.connection.AccountsMgmt().V1().Subscriptions().List()
	clusterInfoList, err := c.executeSubscriptionListRequest(subscriptionListRequest, searchQuery)
	if err != nil {
		c.logger.Errorf("GetSingleClusterInfoForOrganization, err=%s, external_cluster_id=%s", err.Error(), clusterID)
		return
	}
	if clusterInfoList == nil {
		return clusterInfo, fmt.Errorf("cluster not found, cluster_id=%s", clusterID)
	}

	for _, v := range clusterInfoList {
		clusterInfo = v
		break
	}

	c.logger.Debugf(
		"GetSingleClusterInfoForOrganization from AMS API took %s, external_cluster_id=%s", time.Since(tStart), clusterID,
	)
	return clusterInfo, nil
}

// GetInternalOrgIDFromExternal will retrieve the internal organization ID from an external one using AMS API
func (c *amsClientImpl) GetInternalOrgIDFromExternal(orgID string) (string, error) {
	c.logger.Debugf(
		"looking for the internal organization ID for an external one, org_id=%s", orgID,
	)
	orgsListRequest := c.connection.AccountsMgmt().V1().Organizations().List()
	response, err := orgsListRequest.
		Search(fmt.Sprintf("external_id = %s", orgID)).
		Fields("id,external_id").
		Send()

	if err != nil {
		c.logger.Errorf("GetInternalOrgIDFromExternal, err=%s, org_id=%s", err.Error(), orgID)
		return "", err
	}

	if response.Items().Len() != 1 {
		c.logger.Errorf("%s, org_id=%s", orgMoreInternalOrgs, orgID)
		return "", fmt.Errorf(orgMoreInternalOrgs)
	}

	internalID, ok := response.Items().Get(0).GetID()
	if !ok {
		c.logger.Errorf("%s, org_id=%s", orgNoInternalID, orgID)
		return "", fmt.Errorf(orgNoInternalID)
	}

	return internalID, nil
}

func (c *amsClientImpl) executeSubscriptionListRequest(
	subscriptionListRequest *accMgmt.SubscriptionsListRequest,
	searchQuery string,
) (
	clusterInfoMap map[string]ClusterInfo,
	err error,
) {
	clusterInfoMap = map[string]ClusterInfo{}
	for pageNum := 1; ; pageNum++ {
		var err error
		subscriptionListRequest = subscriptionListRequest.
			Size(c.pageSize).
			Page(pageNum).
			Search(searchQuery)

		response, err := subscriptionListRequest.Send()

		if err != nil {
			return clusterInfoMap, err
		}

		// When an empty page is returned, then exit the loop
		if response.Size() == 0 {
			break
		}

		for _, item := range response.Items().Slice() {
			clusterInfo := ClusterInfo{}
			clusterIDstr, ok := item.GetExternalClusterID()
			// we could exclude empty external_cluster_id in the query, but we want to log these special clusters
			if !ok || clusterIDstr == "" {
				if id, ok := item.GetID(); ok {
					c.logger.Warnf("cluster has no external ID, internal_id=%s", id)
				} else {
					c.logger.Errorf("no external or internal cluster ID. cluster=[%v]", item)
				}
				continue
			}
			clusterInfo.ID = clusterIDstr

			displayName, ok := item.GetDisplayName()
			if !ok {
				displayName = string(clusterIDstr)
			}
			clusterInfo.DisplayName = displayName

			status, ok := item.GetStatus()
			if ok {
				clusterInfo.Status = status
			} else {
				c.logger.Warnf("cannot retrieve status of cluster, cluster_id=%s", clusterIDstr)
			}

			plan, ok := item.GetPlan()
			if !ok {
				c.logger.Warnf("cannot retrieve plan of cluster, cluster_id=%s", clusterIDstr)
			}

			if plan != nil {
				clusterType, ok := plan.GetType()
				if ok {
					clusterInfo.Type = clusterType
				} else {
					c.logger.Warnf("cannot retrieve type of cluster, cluster_id=%s", clusterIDstr)
				}
			}

			metrics, ok := item.GetMetrics()
			if ok && len(metrics) > 0 {
				metrics := metrics[0]
				version, ok := metrics.GetOpenshiftVersion()
				if ok {
					clusterInfo.Version = version
				} else {
					c.logger.Warnf("cannot retrieve version of cluster, cluster_id=%s", clusterIDstr)
				}
				provider, ok := metrics.GetCloudProvider()
				if ok {
					clusterInfo.Provider = provider
				} else {
					c.logger.Warnf("cannot retrieve cloud provider of cluster, cluster_id=%s", clusterIDstr)
				}
				region, ok := metrics.GetRegion()
				if ok {
					clusterInfo.Region = region
				} else {
					c.logger.Warnf("cannot retrieve region of cluster, cluster_id=%s", clusterIDstr)
				}
			}

			clusterInfoMap[clusterIDstr] = clusterInfo
		}
	}
	return
}
