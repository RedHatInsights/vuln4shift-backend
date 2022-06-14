package base

import (
	"app/base/models"
	"errors"
	"strconv"
	"strings"
	"time"
)

// ParseBoolArray parses bool array in query arguments,
// array can be checked for max values by limit
func ParseBoolArray(rawValues []string, limit *int) ([]bool, error) {
	if limit != nil && len(rawValues) != *limit {
		return []bool{}, errors.New("invalid bool array format")
	}
	var res []bool
	for _, rawVal := range rawValues {
		val, err := strconv.ParseBool(rawVal)
		if err != nil {
			return res, errors.New("invalid bool value in bool array")
		}
		res = append(res, val)
	}
	return res, nil
}

// ParseDateRange parses 2 member array with date range
func ParseDateRange(rawValues []string) ([]time.Time, error) {
	if len(rawValues) != 2 {
		return []time.Time{}, errors.New("invalid date range format")
	}

	var dateFrom time.Time
	if rawValues[0] == "" {
		dateFrom = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	} else {
		d, err := time.Parse(DateFormat, rawValues[0])
		if err != nil {
			return []time.Time{}, errors.New("invalid date format")
		}
		dateFrom = d
	}

	var dateTo time.Time
	if rawValues[1] == "" {
		dateTo = time.Date(2070, time.January, 1, 0, 0, 0, 0, time.UTC)
	} else {
		d, err := time.Parse(DateFormat, rawValues[1])
		if err != nil {
			return []time.Time{}, errors.New("invalid date format")
		}
		dateTo = d
	}

	return []time.Time{dateFrom, dateTo}, nil
}

// ParseSeverity parses array of severity strings
func ParseSeverity(rawValues []string) ([]models.Severity, error) {
	var res []models.Severity
	for _, raw := range rawValues {
		raw = strings.ToLower(raw)
		switch raw {
		case "", "null":
			res = append(res, models.NotSet)
		case "none":
			res = append(res, models.None)
		case "low":
			res = append(res, models.Low)
		case "medium":
			res = append(res, models.Medium)
		case "moderate":
			res = append(res, models.Moderate)
		case "important":
			res = append(res, models.Important)
		case "high":
			res = append(res, models.High)
		case "critical":
			res = append(res, models.Critical)
		default:
			return res, errors.New("invalid severity argument")
		}
	}
	return res, nil
}

// ParseClusterSeverity parses array of cluster severity strings
func ParseClusterSeverity(rawValues []string) ([]models.Severity, error) {
	var res []models.Severity
	for _, raw := range rawValues {
		raw = strings.ToLower(raw)
		switch raw {
		case "low":
			res = append(res, models.Low)
		case "moderate":
			res = append(res, models.Moderate)
		case "important":
			res = append(res, models.Important)
		case "critical":
			res = append(res, models.Critical)
		default:
			return res, errors.New("invalid cluster_severity argument")
		}
	}
	return res, nil
}

// ParseCvssScoreRange parses array of two member range of cvss score floats
func ParseCvssScoreRange(rawValues []string) ([]float32, error) {
	if len(rawValues) != 2 {
		return []float32{}, errors.New("invalid cvss_score range format")
	}

	var scoreFrom float32
	if rawValues[0] == "" {
		scoreFrom = 0.0
	} else {
		f, err := strconv.ParseFloat(rawValues[0], 32)
		if err != nil {
			return []float32{}, errors.New("invalid cvss score from value")
		}
		scoreFrom = float32(f)
	}

	var scoreTo float32
	if rawValues[1] == "" {
		scoreTo = 10.0
	} else {
		f, err := strconv.ParseFloat(rawValues[1], 32)
		if err != nil {
			return []float32{}, errors.New("invalid cvss score from value")
		}
		scoreTo = float32(f)
	}

	return []float32{float32(scoreFrom), float32(scoreTo)}, nil
}

// ParseUint parses string to int64
func ParseUint(rawValues []string) (uint64, error) {
	if len(rawValues) != 1 {
		return 0, errors.New("Invalid integer")
	}
	res, err := strconv.ParseUint(rawValues[0], 10, 64)
	return res, err
}

// ParseSortArray parses sort params
// +column -> order by column asc
// -column / column -> order by column desc
func ParseSortArray(rawValues []string) []SortItem {
	var res []SortItem
	for _, raw := range rawValues {
		if len(raw) > 0 {
			var item SortItem
			if strings.HasPrefix(raw, "-") {
				item = SortItem{Column: raw[1:], Desc: true}
			} else {
				item = SortItem{Column: raw, Desc: false}
			}
			res = append(res, item)
		}
	}
	return res
}

// ErrInvalidFilterArgument represents error when invalid argument is recieved
var ErrInvalidFilterArgument = errors.New("invalid filter argument")

// ParseFilter parses query argument with name rawName and with rawValues
// filter=1,2,3&filter=cve -> rawName="filter" , rawValues=["1", "2", "3", "cve"]
func ParseFilter(rawName string, rawValues []string) (Filter, error) {
	raw := strings.ToLower(rawName)
	switch raw {
	case SearchQuery:
		if len(rawValues) != 1 {
			return &Search{}, errors.New("invalid search parameter")
		}
		return &Search{RawFilter{raw, rawValues}, rawValues[0]}, nil
	case PublishedQuery:
		dateRange, err := ParseDateRange(rawValues)
		if err != nil {
			return &CvePublishDate{}, errors.New("invalid published parameter format")
		}
		return &CvePublishDate{RawFilter{raw, rawValues}, dateRange[0], dateRange[1]}, nil
	case SeverityQuery:
		severities, err := ParseSeverity(rawValues)
		if err != nil {
			return &Severity{}, err
		}
		return &Severity{RawFilter{raw, rawValues}, severities}, nil
	case ClusterSeverityQuery:
		severities, err := ParseClusterSeverity(rawValues)
		if err != nil {
			return &ClusterSeverity{}, err
		}
		return &ClusterSeverity{RawFilter{raw, rawValues}, severities}, nil
	case CvssScoreQuery:
		scoreRange, err := ParseCvssScoreRange(rawValues)
		if err != nil {
			return &CvssScore{}, err
		}
		return &CvssScore{RawFilter{raw, rawValues}, scoreRange[0], scoreRange[1]}, nil
	case AffectedClustersQuery:
		arrLen := 2
		boolArr, err := ParseBoolArray(rawValues, &arrLen)
		if err != nil {
			return &AffectingClusters{}, errors.New("invalid affected_clusters bool parameter")
		}
		return &AffectingClusters{RawFilter{raw, rawValues}, boolArr[0], boolArr[1]}, nil
	case AffectedImagesQuery:
		arrLen := 2
		boolArr, err := ParseBoolArray(rawValues, &arrLen)
		if err != nil {
			return &AffectingImages{}, errors.New("invalid affected_images bool parameter")
		}
		return &AffectingImages{RawFilter{raw, rawValues}, boolArr[0], boolArr[1]}, nil
	case LimitQuery:
		limit, err := ParseUint(rawValues)
		if err != nil {
			return &Limit{}, errors.New("invalid limit parameter")
		}
		return &Limit{RawFilter{raw, rawValues}, limit}, nil
	case OffsetQuery:
		offset, err := ParseUint(rawValues)
		if err != nil {
			return &Offset{}, errors.New("invalid offset parameter")
		}
		return &Offset{RawFilter{raw, rawValues}, offset}, nil
	case SortQuery:
		sortArr := ParseSortArray(rawValues)
		return &Sort{RawFilter{raw, rawValues}, sortArr}, nil
	default:
		return &Search{}, ErrInvalidFilterArgument
	}
}
