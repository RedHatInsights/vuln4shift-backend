package pyxis

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"app/base/logging"
	"app/base/models"
	"app/base/utils"
)

var (
	logger *logrus.Logger
)

func init() {
	var err error
	logger, err = logging.CreateLogger(utils.Cfg.LoggingLevel)
	if err != nil {
		fmt.Println("Error setting up logger.")
		os.Exit(1)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func registerMissingCves(tx *gorm.DB, apiImageCves []string) error {
	toInsertCves := []models.Cve{}
	var found bool
	for _, cveName := range apiImageCves {
		if _, found = dbCveMap[cveName]; !found {
			if _, found = dbCveMapPending[cveName]; !found {
				toInsertCves = append(toInsertCves, models.Cve{Name: cveName, Description: "unknown", Severity: models.NotSet, Cvss2Score: 0.0, Cvss3Score: 0.0})
			}
		}
	}

	toInsertCvesCnt := len(toInsertCves)
	logger.Debugf("CVEs to insert: %d", toInsertCvesCnt)

	if toInsertCvesCnt > 0 {
		// Use conflict clause as the cve table can be changed from vmaas-sync
		// TODO: needs to be sorted insert to avoid deadlocks
		if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "name"}}, DoNothing: true}).Create(&toInsertCves).Error; err != nil {
			syncError.WithLabelValues(dbRegisterMissingCves).Inc()
			return err
		}
		missingCvesRegistered.Add(float64(toInsertCvesCnt))

	}

	// Add newly inserted CVEs to the cache after commit
	for _, cve := range toInsertCves {
		dbCveMapPending[cve.Name] = cve
	}

	return nil
}

func syncImage(tx *gorm.DB, image models.Image) error {
	if image.ID == 0 {
		if err := tx.Create(&image).Error; err != nil {
			syncError.WithLabelValues(dbInsert).Inc()
			return err
		}
		dbImageMapPending[image.PyxisID] = image // Add newly inserted image to the cache after commit
	} else {
		if err := tx.Save(&image).Error; err != nil {
			syncError.WithLabelValues(dbUpdate).Inc()
			return err
		}
	}

	if utils.Cfg.SkipImageCveSync {
		return nil
	}

	apiImageCves, err := getAPIImageCves(image.PyxisID)
	if err != nil {
		return err
	}

	if err := registerMissingCves(tx, apiImageCves); err != nil {
		return err
	}

	dbImageCveMap, err := getDbImageCves(image.ID)
	if err != nil {
		return err
	}

	toInsertImageCves := []models.ImageCve{}
	toDeleteImageCves := []models.ImageCve{}
	var cve models.Cve
	var found bool
	for _, cveName := range apiImageCves {
		// Lookup CVE in the cache (also in the pending cache)
		if cve, found = dbCveMap[cveName]; !found {
			if cve, found = dbCveMapPending[cveName]; !found {
				syncError.WithLabelValues(dbCveNotInCache).Inc()
				err := fmt.Errorf("CVE not in cache: %s", cveName)
				return err
			}
		}
		if _, found := dbImageCveMap[cve.ID]; !found {
			// image_cve pair not found
			toInsertImageCves = append(
				toInsertImageCves,
				models.ImageCve{
					ImageID: image.ID,
					CveID:   cve.ID,
				},
			)
		} else {
			delete(dbImageCveMap, cve.ID)
		}
	}

	// Delete the rest of image_cve pairs in DB not returned by API
	for _, imageCve := range dbImageCveMap {
		toDeleteImageCves = append(toDeleteImageCves, imageCve)
	}

	toInsertImageCvesCnt := len(toInsertImageCves)
	toDeleteImageCvesCnt := len(toDeleteImageCves)

	logger.Debugf("Image-CVE pairs to insert: %d", toInsertImageCvesCnt)
	logger.Debugf("Image-CVE pairs to delete: %d", toDeleteImageCvesCnt)

	if toInsertImageCvesCnt > 0 {
		if err := tx.Create(&toInsertImageCves).Error; err != nil {
			syncError.WithLabelValues(dbInsert).Inc()
			return err
		}
		imageCvesInserted.Add(float64(toInsertImageCvesCnt))
	}

	if toDeleteImageCvesCnt > 0 {
		deleteTx := tx.Session(&gorm.Session{})
		for _, image := range toDeleteImageCves {
			deleteTx = deleteTx.Or(&image)
		}
		if err := deleteTx.Delete(&models.ImageCve{}).Error; err != nil {
			syncError.WithLabelValues(dbDelete).Inc()
			return err
		}
		imageCvesDeleted.Add(float64(toDeleteImageCvesCnt))
	}

	return nil
}

func syncRepo(repo models.Repository) error {
	// Repository is our database unit, commit once per every repo
	tx := DB.Begin()
	// Do a rollback by default (don't need to specify on every return), will do nothing when everything is committed
	defer tx.Rollback()

	if repo.ID == 0 {
		if err := tx.Create(&repo).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Save(&repo).Error; err != nil {
			return err
		}
	}

	apiRepoImages, err := getAPIRepoImages(repo.Registry, repo.Repository)
	if err != nil {
		return err
	}

	toSyncImages := []models.Image{}

	var dbArch models.Arch
	var found bool
	for pyxisID, apiImage := range apiRepoImages {
		if dbArch, found = dbArchMap[apiImage.Architecture]; !found {
			if dbArch, found = dbArchMapPending[apiImage.Architecture]; !found {
				if len(apiImage.Architecture) == 0 {
					err := fmt.Errorf("empty arch for image, pyxis_id=%s", apiImage.PyxisID)
					return err
				}
				dbArch = models.Arch{Name: apiImage.Architecture}
				if err := tx.Create(&dbArch).Error; err != nil {
					return err
				}
				dbArchMapPending[apiImage.Architecture] = dbArch
			}
		}
		manifestListDigest := apiImage.Repositories[0].ManifestListDigest
		manifestSchema2Digest := apiImage.Repositories[0].ManifestSchema2Digest
		dockerImageDigest := apiImage.DockerImageDigest
		var manifestListDigestDB, manifestSchema2DigestDB, dockerImageDigestDB *string
		if len(manifestListDigest) > 0 {
			manifestListDigestDB = &manifestListDigest
		}
		if len(manifestSchema2Digest) > 0 {
			manifestSchema2DigestDB = &manifestSchema2Digest
		}
		if len(dockerImageDigest) > 0 {
			dockerImageDigestDB = &dockerImageDigest
		}
		modifiedDate := apiImage.ModifiedDate
		if dbImage, found := dbImageMap[pyxisID]; !found {
			if utils.Cfg.SkipImageCveSync {
				modifiedDate = time.Time{} // Set to zero value if not syncing CVEs
			}
			toSyncImages = append(
				toSyncImages,
				models.Image{
					PyxisID:               apiImage.PyxisID,
					ModifiedDate:          modifiedDate,
					ManifestListDigest:    manifestListDigestDB,
					ManifestSchema2Digest: manifestSchema2DigestDB,
					DockerImageDigest:     dockerImageDigestDB,
					ArchID:                dbArch.ID,
				},
			)
		} else if apiImage.ModifiedDate.After(dbImage.ModifiedDate) || utils.Cfg.ForceSync {
			if utils.Cfg.SkipImageCveSync {
				modifiedDate = dbImage.ModifiedDate // Keep the original value if not syncing CVEs
			}
			dbImage.ModifiedDate = modifiedDate
			dbImage.ManifestListDigest = manifestListDigestDB
			dbImage.ManifestSchema2Digest = manifestSchema2DigestDB
			dbImage.DockerImageDigest = dockerImageDigestDB
			dbImage.ArchID = dbArch.ID
			toSyncImages = append(toSyncImages, dbImage)
		}
	}

	logger.Debugf("Images to sync: %d", len(toSyncImages))

	for _, image := range toSyncImages {
		err := syncImage(tx, image)
		if err != nil {
			return err
		}
	}

	// Sync also Repository - Image pairs
	dbRepositoryImageMap, err := getDbRepositoryImages(repo.ID)
	if err != nil {
		return err
	}

	toInsertRepositoryImages := []models.RepositoryImage{}
	toDeleteRepositoryImages := []models.RepositoryImage{}

	var image models.Image
	for pyxisID := range apiRepoImages {
		// Lookup image in the cache (also in the pending cache because it might be inserted in current transaction)
		if image, found = dbImageMap[pyxisID]; !found {
			if image, found = dbImageMapPending[pyxisID]; !found {
				syncError.WithLabelValues(dbImageNotInCache).Inc()
				err := fmt.Errorf("image not in cache: %s", pyxisID)
				return err
			}
		}
		if _, found := dbRepositoryImageMap[image.ID]; !found {
			// repository_image pair not found
			toInsertRepositoryImages = append(
				toInsertRepositoryImages,
				models.RepositoryImage{
					RepositoryID: repo.ID,
					ImageID:      image.ID,
				},
			)
		} else {
			delete(dbRepositoryImageMap, image.ID)
		}
	}

	// Delete the rest of repository_image pairs in DB not returned by API
	for _, repositoryImage := range dbRepositoryImageMap {
		toDeleteRepositoryImages = append(toDeleteRepositoryImages, repositoryImage)
	}

	toInsertRepositoryImagesCnt := len(toInsertRepositoryImages)
	toDeleteRepositoryImagesCnt := len(toDeleteRepositoryImages)

	logger.Debugf("Repository-Image pairs to insert: %d", toInsertRepositoryImagesCnt)
	logger.Debugf("Repository-Image pairs to delete: %d", toDeleteRepositoryImagesCnt)

	if toInsertRepositoryImagesCnt > 0 {
		if err := tx.Create(&toInsertRepositoryImages).Error; err != nil {
			syncError.WithLabelValues(dbInsert).Inc()
			return err
		}
		syncedImages.WithLabelValues(repo.Repository).Add(float64(toInsertRepositoryImagesCnt))
	}

	if toDeleteRepositoryImagesCnt > 0 {
		deleteTx := tx.Session(&gorm.Session{})
		for _, repo := range toDeleteRepositoryImages {
			deleteTx = deleteTx.Or(&repo)
		}
		if err := deleteTx.Delete(&models.RepositoryImage{}).Error; err != nil {
			syncError.WithLabelValues(dbDelete).Inc()
			return err
		}
		deletedImages.WithLabelValues(repo.Repository).Add(float64(toDeleteRepositoryImagesCnt))
	}

	return tx.Commit().Error
}

func syncRepos() {
	apiRepoMap, err := getAPIRepositories()
	if err != nil {
		logger.Fatalf("Unable to get repositories from Pyxis: %s", err)
	}
	logger.Infof("Repositories in Pyxis: %d", len(apiRepoMap))

	toSyncRepos := []models.Repository{}

	for repoKey, apiRepo := range apiRepoMap {
		if passed := repositoryInProfile(apiRepo.Registry, apiRepo.Repository); !passed {
			continue
		} else if dbRepo, found := dbRepoMap[repoKey]; !found {
			// repo is not found by its registry/repo key, but through pyxis_id, update it
			if dbRepo, found := dbPyxisIDRepoMap[apiRepo.PyxisID]; found {
				dbRepo.PyxisID = apiRepo.PyxisID
				dbRepo.ModifiedDate = apiRepo.ModifiedDate
				dbRepo.Registry = apiRepo.Registry
				dbRepo.Repository = apiRepo.Repository
				toSyncRepos = append(toSyncRepos, dbRepo)
			} else {
				// completely new repo, insert it
				toSyncRepos = append(
					toSyncRepos,
					models.Repository{
						PyxisID:      apiRepo.PyxisID,
						ModifiedDate: apiRepo.ModifiedDate,
						Registry:     apiRepo.Registry,
						Repository:   apiRepo.Repository,
					},
				)
			}
		} else if dbRepo, found := dbPyxisIDRepoMap[dbRepo.PyxisID]; !found {
			// repo is found by its registry/repo key, but not through pyxis_id, update it
			dbRepo.PyxisID = apiRepo.PyxisID
			dbRepo.ModifiedDate = apiRepo.ModifiedDate
			dbRepo.Registry = apiRepo.Registry
			dbRepo.Repository = apiRepo.Repository
			toSyncRepos = append(toSyncRepos, dbRepo)
		} else if apiRepo.ModifiedDate.After(dbRepo.ModifiedDate) || utils.Cfg.ForceSync {
			// repo is found by its registry/repo key and pyxis_id but got modified, update it
			dbRepo.PyxisID = apiRepo.PyxisID
			dbRepo.ModifiedDate = apiRepo.ModifiedDate
			dbRepo.Registry = apiRepo.Registry
			dbRepo.Repository = apiRepo.Repository
			toSyncRepos = append(toSyncRepos, dbRepo)
			delete(dbRepoMap, repoKey)
		} else {
			// repo is not found by its registry/repo key and pyxis_id, keep it untouched
			delete(dbRepoMap, repoKey)
		}
	}

	toSyncReposCnt := len(toSyncRepos)
	logger.Infof("Repositories to sync (profile=%s): %d", profile, toSyncReposCnt)
	logger.Infof("Repositories in DB not known to Pyxis or not in current profile (profile=%s): %d", profile, len(dbRepoMap))

	for i, repo := range toSyncRepos {
		logger.Infof("Syncing repo: repo=%s/%s [%d/%d]", repo.Registry, repo.Repository, i+1, toSyncReposCnt)
		if err := syncRepo(repo); err != nil {
			logger.Warnf("Syncing repo failed, skipping: repo=%s/%s, err=%s", repo.Registry, repo.Repository, err)
			emptyPendingCache() // Not successfully committed, don't update cache
		} else {
			flushPendingCache() // Update cache
		}
	}
}

func Start() {
	logger.Info("Starting Pyxis sync.")

	pusher := GetMetricsPusher()

	parseProfiles() // Parse static yaml file with profiles (list of repositories)

	if err := dbConfigure(); err != nil {
		syncError.WithLabelValues(dbConnection).Inc()
		logger.Fatalf("Unable to get GORM connection: %s", err)
	}
	if err := prepareMaps(); err != nil {
		syncError.WithLabelValues(dbFetch).Inc()
		logger.Fatalf("Unable to fetch data from DB: %s", err)
	}

	logger.Infof("Repositories in DB: %d", len(dbRepoMap))
	logger.Infof("Images in DB: %d", len(dbImageMap))
	logger.Infof("CVEs in DB: %d", len(dbCveMap))

	syncRepos()

	logger.Info("Finished Pyxis sync.")

	if err := pusher.Add(); err != nil {
		logger.Warningln("Could not push to Pushgateway:", err)
	}
}
