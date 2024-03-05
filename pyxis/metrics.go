package pyxis

import (
	"app/base/utils"

	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	dbConnection              = "db-connection"
	dbFetch                   = "db-fetch"
	dbInsert                  = "db-insert"
	dbUpdate                  = "db-update"
	dbDelete                  = "db-delete"
	dbImageCveNotFound        = "db-image-cve-not-found"
	dbRepositoryImageNotFound = "db-repository-image-not-found"
	dbImageNotInCache         = "db-image-not-in-cache"
	dbCveNotInCache           = "db-cve-not-in-cache"
	dbRegisterMissingCves     = "db-register-missing-cves"

	job = "pyxis"
)

var (
	syncError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many errors and of which type",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "sync_error",
	}, []string{"type"})

	pyxisRequestError = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "Pyxis api request error by status code",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "request_error",
	}, []string{"url", "method", "code"})

	syncedImages = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many images were synced",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "sync_images",
	}, []string{"repo"})

	deletedImages = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many images were deleted",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "delete_images",
	}, []string{"repo"})

	updatedImages = prometheus.NewCounterVec(prometheus.CounterOpts{
		Help:      "How many images were updated",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "update_images",
	}, []string{"repo"})

	missingCvesRegistered = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many missing cves were registered",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "missing_cves_registered",
	})

	imageCvesDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many image cves pairs deleted",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "images_cves_deleted",
	})

	imageCvesInserted = prometheus.NewCounter(prometheus.CounterOpts{
		Help:      "How many image cves pairs inserted",
		Namespace: "vuln4shift",
		Subsystem: "pyxis",
		Name:      "images_cves_inserted",
	})
)

func getMetricsPusher() *push.Pusher {
	return utils.GetMetricsPusher(
		job,
		syncError,
		pyxisRequestError,
		syncedImages,
		deletedImages,
		updatedImages,
		missingCvesRegistered,
		imageCvesDeleted,
		imageCvesInserted)
}
