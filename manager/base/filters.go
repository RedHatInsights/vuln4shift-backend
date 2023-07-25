package base

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"app/base/models"

	"gorm.io/gorm"
)

const DateFormat = "2006-01-02"
const (
	SearchQuery           = "search"
	PublishedQuery        = "published"
	SeverityQuery         = "severity"
	ClusterSeverityQuery  = "cluster_severity"
	CvssScoreQuery        = "cvss_score"
	AffectedClustersQuery = "affected_clusters"
	LimitQuery            = "limit"
	OffsetQuery           = "offset"
	SortQuery             = "sort"
	DataFormatQuery       = "data_format"
	ReportQuery           = "report"
	ProviderQuery         = "provider"
	StatusQuery           = "status"
	VersionQuery          = "version"
	ExploitsQuery         = "exploits"
)

const (
	SortFilterArgs = "sort_filter"
)

const (
	CveSearch             = "CveSearch"
	ExposedClustersSearch = "ExposedClustersSearch"
)

const (
	JSONFormat = iota
	CSVFormat
)

// Filter interface, represents filter obtained from
// query argument in request link.
type Filter interface {
	ApplyQuery(*gorm.DB, map[string]interface{}) error
	RawQueryVal() string
	RawQueryVals() []string
	RawQueryName() string
}

// RawFilter implements Filter interface, contains
// raw name of query argument and raw parsed query
// values in string
type RawFilter struct {
	RawParam  string
	RawValues []string
}

// RawQueryName getter to the parameter name in query
func (b *RawFilter) RawQueryName() string {
	return b.RawParam
}

// RawQueryVal returns obtained raw values formatted in query value string
func (b *RawFilter) RawQueryVal() string {
	return strings.Join(b.RawValues[:], ",")
}

// RawQueryVals returns parsed raw values from query value string
func (b *RawFilter) RawQueryVals() []string {
	return b.RawValues
}

// Search represents filter for CVE substring search
// ex. search=CVE-2022
type Search struct {
	RawFilter
	value string
}

// ApplyQuery filters CVEs by their substring match name or description
func (c *Search) ApplyQuery(tx *gorm.DB, args map[string]interface{}) error {
	regex := fmt.Sprintf("%%%s%%", c.value)

	switch args[SearchQuery] {
	case CveSearch:
		tx.Where("cve.name ILIKE ? OR cve.description ILIKE ?", regex, regex)
		return nil
	case ExposedClustersSearch:
		tx.Where("cluster.display_name ILIKE ? OR cluster.uuid::varchar ILIKE ?", regex, regex)
		return nil
	}
	return nil
}

// CvePublishDate represents filter for CVE publish date filtering
// ex: published=2021-01-01,2022-02-02
type CvePublishDate struct {
	RawFilter
	From time.Time
	To   time.Time
}

// ApplyQuery filters CVEs by their public date limit
func (c *CvePublishDate) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	tx.Where("cve.public_date >= ? AND cve.public_date <= ?", c.From, c.To)
	return nil
}

// Severity represents CVE severity filter
// ex. severity=critical,important,none
type Severity struct {
	RawFilter
	Value []models.Severity
}

// ApplyQuery filters CVEs by their severity
func (s *Severity) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	tx.Where("cve.severity IN ?", s.Value)
	return nil
}

// ClusterSeverity represents CVE severity filter for clusters
// ex. cluster_severity=critical,important
type ClusterSeverity struct {
	RawFilter
	Value []models.Severity
}

func (s *ClusterSeverity) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	severities := []string{}
	for _, severity := range s.Value {
		switch severity {
		case models.Critical:
			severities = append(severities, "cluster.cve_cache_critical > 0")
		case models.Important:
			severities = append(severities, "cluster.cve_cache_important > 0")
		case models.Moderate:
			severities = append(severities, "cluster.cve_cache_moderate > 0")
		case models.Low:
			severities = append(severities, "cluster.cve_cache_low > 0")
		}
	}
	if len(severities) > 0 {
		tx.Where(strings.Join(severities, " OR "))
	}
	return nil
}

// CvssScore represents filter for CVE cvss2/3 score range
// cvss_score=0.0,9.0
type CvssScore struct {
	RawFilter
	From float32
	To   float32
}

// ApplyQuery filters CVEs by cvss2/3 score range
func (c *CvssScore) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	tx.Where("GREATEST(cve.cvss3_score, cve.cvss2_score) >= ? AND GREATEST(cve.cvss3_score, cve.cvss2_score) <= ?", c.From, c.To)
	return nil
}

// AffectingClusters represents filter for count of affected clusters
// ex. clusters_exposed=true,true
type AffectingClusters struct {
	RawFilter
	OneOrMore bool
	None      bool
}

// ApplyQuery filters rows by count of affected clusters
func (a *AffectingClusters) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	if a.None == a.OneOrMore {
		return nil
	}

	if a.None {
		tx.Where("COALESCE(ce, 0) = 0")
	} else {
		tx.Where("COALESCE(ce, 0) != 0")
	}

	if a.OneOrMore {
		tx.Where("COALESCE(ce, 0) > 0")
	} else {
		tx.Where("COALESCE(ce, 0) = 0")
	}

	return nil
}

// Limit filter sets number of data objects per page
// ex. limit=20
type Limit struct {
	RawFilter
	Value uint64
}

// ApplyQuery limits the number of data in query - limit per page
func (l *Limit) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	tx.Limit(int(l.Value))
	return nil
}

// Offset filter sets an offset of data in query - start of the page
// ex. offset=40
type Offset struct {
	RawFilter
	Value uint64
}

// ApplyQuery sets and offset from the rows result
func (o *Offset) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	tx.Offset(int(o.Value))
	return nil
}

// SortItem represents an single column row sort expression
// Used by the Sort filter
type SortItem struct {
	Column string
	Desc   bool
}

// SortArgs represents an argument for Sort filter
// SortableColumns represents mapping from user selected column
// to the correct sql expression column
// DefaultSortable contains a default sorting defined by controller
type SortArgs struct {
	SortableColumns map[string]string
	DefaultSortable []SortItem
}

// Sort filter sorts a query by given list of sort item expressions
// ex. sort=synopsis,cvss_score
type Sort struct {
	RawFilter
	Values []SortItem
}

// ApplyQuery sorts the resulting query, query is sorted
// 1st - by user defined columns
// 2nd - by controller selected default columns
func (s *Sort) ApplyQuery(tx *gorm.DB, args map[string]interface{}) error {
	if i, exists := args[SortFilterArgs]; exists {
		sortArgs, ok := i.(SortArgs)
		if !ok {
			return nil
		}
		// Sort by user selected columns
		for _, item := range s.Values {
			// Check if selected user column is mappable to sortable column sql expression
			if col, exists := sortArgs.SortableColumns[item.Column]; exists {
				if item.Desc {
					tx.Order(fmt.Sprintf("%s DESC NULLS LAST", col))
				} else {
					tx.Order(fmt.Sprintf("%s ASC NULLS LAST", col))
				}
			} else {
				return errors.New("invalid sort column selected")
			}
		}
		// Sort by default sortable
		for _, item := range sortArgs.DefaultSortable {
			if col, exists := sortArgs.SortableColumns[item.Column]; exists {
				if item.Desc {
					tx.Order(fmt.Sprintf("%s DESC NULLS LAST", col))
				} else {
					tx.Order(fmt.Sprintf("%s ASC NULLS LAST", col))
				}
			}
		}
	}
	return nil
}

// DataFormat filter represents serialisation format of data section in response
// ex. data_format=csv
type DataFormat struct {
	RawFilter
	Value uint64
}

// ApplyQuery format does not do anything to the DB query
func (d *DataFormat) ApplyQuery(_ *gorm.DB, _ map[string]interface{}) error {
	return nil
}

// Report filter represents ignoring of all paging query to return all
// for json/csv report export
type Report struct {
	RawFilter
	Value bool
}

// ApplyQuery resets limit and offset on the query
func (r *Report) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	if r.Value {
		tx.Limit(-1)
		tx.Offset(-1)
	}
	return nil
}

// Provider represents cluster provider filter
type Provider struct {
	RawFilter
	Values []string
}

// ApplyQuery filters clusters by provider
func (p *Provider) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	if len(p.Values) > 0 {
		tx.Where("COALESCE(cluster.provider, 'N/A') IN ?", p.Values)
	}
	return nil
}

// Status represents cluster status filter
// ex. status=disconnected,connected
// FIXME: values should be enums of the possible cluster status values
type Status struct {
	RawFilter
	Values []string
}

// ApplyQuery filters the cluster by its status
func (s *Status) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	if len(s.Values) > 0 {
		tx.Where("cluster.status IN ?", s.Values)
	}
	return nil
}

// Version filter represents version of cluster
// ex. version=4.1,4.2
type Version struct {
	RawFilter
	Values []string
}

// ApplyQuery searches for version in array
func (v *Version) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	if len(v.Values) > 0 {
		tx.Where("cluster.version IN ?", v.Values)
	}
	return nil
}

// Exploits filters only CVEs with known exploits
// ex. exploits=true
type Exploits struct {
	RawFilter
	Value bool
}

// ApplyQuery searches for version in an array
func (e *Exploits) ApplyQuery(tx *gorm.DB, _ map[string]interface{}) error {
	if e.Value {
		tx.Where("cve.exploit_data IS NOT NULL")
	} else {
		tx.Where("cve.exploit_data IS NULL")
	}
	return nil
}

// ApplyFilters applies requested filters from query params on created query from controller,
// filters needs to be allowed from controller in allowedFilters array
// Returns used filters for querying
func ApplyFilters(query *gorm.DB, allowedFilters []string,
	requestedFilters map[string]Filter, args map[string]interface{}) (map[string]Filter, error) {
	usedFilters := make(map[string]Filter)
	for _, allowedFilter := range allowedFilters {
		if filter, requested := requestedFilters[allowedFilter]; requested {
			err := filter.ApplyQuery(query, args)
			if err != nil {
				return map[string]Filter{}, err
			}
			usedFilters[allowedFilter] = filter
		}
	}
	return usedFilters, nil
}
