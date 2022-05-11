package middlewares

import (
	"app/base/models"
	"app/manager/base"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var filterer = Filterer()

func TestFiltererEmpty(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)
	assert.Equal(t, 3, len(filters), "Should be 3 filters, the default ones - sort, limit, offset")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test", nil)

	filterer(ctx)

	filters = base.GetRequestedFilters(ctx)
	assert.Equal(t, 3, len(filters), "Should be 3 filters, the default ones - sort, limit, offset")
}

func TestFiltererInvalidParam(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?nonexisting=param&another=one", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Unknown arguments are caught as error")
}

func TestFiltererValidSearch(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?search=CVE-2022", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)
	f, e := filters[base.SearchQuery]
	assert.Equal(t, e, true, "Should be cve search filter")
	filter, ok := f.(*base.Search)

	assert.Equal(t, true, ok, "Should be cve search filter")
	assert.Equal(t, filter.RawValues, []string{"CVE-2022"}, "Filter should contain CVE-2022 infix")
	assert.Equal(t, 4, len(filters), "Should be 4 filters, the default ones - sort, limit, offset, search one")
}

func TestFiltererInvalidSearch(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?search=CVE-2022,b,d", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Multiple search args, expecting HTTP 400 response")
}

func TestFiltererValidPublished(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?published=1980-01-01,2022-01-01", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)

	f, e := filters[base.PublishedQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok := f.(*base.CvePublishDate)
	assert.Equal(t, true, ok, "Should be cve search filter")

	fromExpected, _ := time.Parse(base.DateFormat, "1980-01-01")
	toExpected, _ := time.Parse(base.DateFormat, "2022-01-01")
	assert.Equal(t, filter.From, fromExpected, "Published should be from date 1980-01-01")
	assert.Equal(t, filter.To, toExpected, "Published should be to date 2022-01-01")
	assert.Equal(t, 4, len(filters), "Should be 4 filters, the default ones - sort, limit, offset, published one")
}

func TestFiltererInvalidPublished(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?published=", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid range, xpecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?published=1,2", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid dates, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?published=1990-01-01,2000-01-01,2020-01-01", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Three dates, expecting HTTP 400 response")
}

func TestFiltererValidSeverity(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?severity=none,low,medium,moderate,important,high,critical", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)

	f, e := filters[base.SeverityQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok := f.(*base.Severity)
	assert.Equal(t, true, ok, "Should be cve search filter")

	assert.Equal(t, filter.Value, []models.Severity{models.None, models.Low, models.Medium, models.Moderate, models.Important, models.High, models.Critical},
		"Expecting 7 severity values")
	assert.Equal(t, 4, len(filters), "Should be 4 filters, the default ones - sort, limit, offset, severity one")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?severity=medium", nil)

	filterer(ctx)

	filters = base.GetRequestedFilters(ctx)

	f, e = filters[base.SeverityQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok = f.(*base.Severity)
	assert.Equal(t, true, ok, "Should be cve search filter")

	assert.Equal(t, filter.Value, []models.Severity{models.Medium},
		"Expecting medium severity value")
	assert.Equal(t, 4, len(filters), "Should be 4 filters, the default ones - sort, limit, offset, severity one")
}

func TestFiltererInvalidSeverity(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?severity=nonexisting", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid severity, xpecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?severity=moderate,nonexisting", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "One valid, one invalid, expecting HTTP 400 response")
}

func TestFiltererValidCvssScore(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?cvss_score=0.5,9.9", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)

	f, e := filters[base.CvssScoreQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok := f.(*base.CvssScore)
	assert.Equal(t, true, ok, "Should be cve search filter")

	fromExpected := float32(0.5)
	toExpected := float32(9.9)
	assert.Equal(t, filter.From, fromExpected, "Expecting cvss_score from 0.5")
	assert.Equal(t, filter.To, toExpected, "Expecting cvss_score to 0.9")
	assert.Equal(t, 4, len(filters), "Should be 4 filters, the default ones - sort, limit, offset, cvss_score one")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?cvss_score=0,10.1", nil)

	filterer(ctx)

	filters = base.GetRequestedFilters(ctx)

	f, e = filters[base.CvssScoreQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok = f.(*base.CvssScore)
	assert.Equal(t, true, ok, "Should be cve search filter")

	fromExpected = float32(0.0)
	toExpected = float32(10.1)
	assert.Equal(t, filter.From, fromExpected, "Expecting cvss_score from 0.5")
	assert.Equal(t, filter.To, toExpected, "Expecting cvss_score to 0.9")
	assert.Equal(t, 4, len(filters), "Should be 4 filters, the default ones - sort, limit, offset, cvss_score one")
}

func TestFiltererInvalidCvssScore(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?cvss_score=", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid range, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?cvss_score=2,B", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid float range, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?cvss_score=3.0,5.0,5.9", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Three values, expecting HTTP 400 response")
}

func TestFiltererValidLimit(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?limit=20", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)

	f, e := filters[base.LimitQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok := f.(*base.Limit)
	assert.Equal(t, true, ok, "Should be cve search filter")

	assert.Equal(t, filter.Value, uint64(20), "Expecting limit to be 20")
	assert.Equal(t, 3, len(filters), "Should be 3 filters, the default ones - sort, limit, offset")
}

func TestFiltererInvalidLimit(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?limit=", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Empty limit, xpecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?limit=A", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid limit, one invalid, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?limit=2,B", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid limit, one invalid, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?limit=-5", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Negative limit, one invalid, expecting HTTP 400 response")
}

func TestFiltererValidOffset(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?offset=40", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)

	f, e := filters[base.OffsetQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok := f.(*base.Offset)
	assert.Equal(t, true, ok, "Should be cve search filter")

	assert.Equal(t, filter.Value, uint64(40), "Expecting limit to be 20")
	assert.Equal(t, 3, len(filters), "Should be 3 filters, the default ones - sort, limit, offset")
}

func TestFiltererInvalidOffset(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?offset=", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Empty limit, xpecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?offset=A", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid limit, one invalid, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?offset=2,B", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Invalid limit, one invalid, expecting HTTP 400 response")

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?offset=-5", nil)

	filterer(ctx)

	base.GetRequestedFilters(ctx)
	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status(), "Negative limit, one invalid, expecting HTTP 400 response")
}

func TestFiltererSort(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?sort=id,-synopsis", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)

	f, e := filters[base.SortQuery]
	assert.Equal(t, e, true, "Should be published filter")
	filter, ok := f.(*base.Sort)
	assert.Equal(t, true, ok, "Should be cve search filter")

	assert.Equal(t, filter.Values, []base.SortItem{{Column: "id", Desc: false}, {Column: "synopsis", Desc: true}}, "Expecting limit to be 20")
	assert.Equal(t, 3, len(filters), "Should be 3 filters, the default ones - sort, limit, offset")
}

func TestFiltererMultiple(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/test?limit=20&offset=40&sort=id&search=CVE&published=1980-01-01,2022-01-01&severity=high,critical&cvss_score=0.0,9.0", nil)

	filterer(ctx)

	filters := base.GetRequestedFilters(ctx)
	assert.Equal(t, 7, len(filters), "Should be 7 filters")
}
