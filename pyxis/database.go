package pyxis

import (
	"fmt"

	"gorm.io/gorm"

	"app/base/models"
	"app/base/utils"
)

var (
	DB         *gorm.DB
	dbRepoMap  map[string]models.Repository
	dbArchMap  map[string]models.Arch
	dbImageMap map[string]models.Image
	dbCveMap   map[string]models.Cve

	// Objects to add to the cached maps above after succesful commit
	dbArchMapPending  = map[string]models.Arch{}
	dbImageMapPending = map[string]models.Image{}
	dbCveMapPending   = map[string]models.Cve{}
)

func dbConfigure() error {
	dsn := utils.GetDbURL(false)
	var err error
	DB, err = models.GetGormConnection(dsn)

	if err != nil {
		return err
	}
	return nil
}

func formatRepoMapKey(registry, repository string) string {
	return fmt.Sprintf("%s/%s", registry, repository)
}

func prepareDbRepositories() error {
	repoRows := []models.Repository{}
	if err := DB.Find(&repoRows).Error; err != nil {
		return err
	}
	dbRepoMap = make(map[string]models.Repository, len(repoRows))
	for _, repo := range repoRows {
		dbRepoMap[formatRepoMapKey(repo.Registry, repo.Repository)] = repo
	}
	return nil
}

func prepareDbArchs() error {
	archRows := []models.Arch{}
	if err := DB.Find(&archRows).Error; err != nil {
		return err
	}
	dbArchMap = make(map[string]models.Arch, len(archRows))
	for _, arch := range archRows {
		dbArchMap[arch.Name] = arch
	}
	return nil
}

func prepareDbImages() error {
	imageRows := []models.Image{}
	if err := DB.Find(&imageRows).Error; err != nil {
		return err
	}
	dbImageMap = make(map[string]models.Image, len(imageRows))
	for _, image := range imageRows {
		dbImageMap[image.PyxisID] = image
	}
	return nil
}

func prepareDbCves() error {
	cveRows := []models.Cve{}
	if err := DB.Order("name").Find(&cveRows).Error; err != nil {
		return err
	}
	dbCveMap = make(map[string]models.Cve, len(cveRows))
	for _, cve := range cveRows {
		dbCveMap[cve.Name] = cve
	}
	return nil
}

func prepareMaps() error {
	if err := prepareDbRepositories(); err != nil {
		return err
	}
	if err := prepareDbArchs(); err != nil {
		return err
	}
	if err := prepareDbImages(); err != nil {
		return err
	}
	if err := prepareDbCves(); err != nil {
		return err
	}
	return nil
}

func emptyPendingCache() {
	dbArchMapPending = map[string]models.Arch{}
	dbImageMapPending = map[string]models.Image{}
	dbCveMapPending = map[string]models.Cve{}
}

func flushPendingCache() {
	for key, val := range dbArchMapPending {
		dbArchMap[key] = val
	}
	for key, val := range dbImageMapPending {
		dbImageMap[key] = val
	}
	for key, val := range dbCveMapPending {
		dbCveMap[key] = val
	}
	emptyPendingCache()
}

func getDbImageCves(imageID int64) (map[int64]models.ImageCve, error) {
	imageCveRows := []models.ImageCve{}
	if err := DB.Where("image_id = ?", imageID).Find(&imageCveRows).Error; err != nil {
		syncError.WithLabelValues(dbImageCveNotFound).Inc()
		return nil, err
	}
	dbImageCveMap := make(map[int64]models.ImageCve, len(imageCveRows))
	for _, imageCve := range imageCveRows {
		dbImageCveMap[imageCve.CveID] = imageCve
	}
	return dbImageCveMap, nil
}

func getDbRepositoryImages(repositoryID int64) (map[int64]models.RepositoryImage, error) {
	repositoryImageRows := []models.RepositoryImage{}
	if err := DB.Where("repository_id = ?", repositoryID).Find(&repositoryImageRows).Error; err != nil {
		syncError.WithLabelValues(dbRepositoryImageNotFound).Inc()
		return nil, err
	}
	dbRepositoryImageMap := make(map[int64]models.RepositoryImage, len(repositoryImageRows))
	for _, repositoryImage := range repositoryImageRows {
		dbRepositoryImageMap[repositoryImage.ImageID] = repositoryImage
	}
	return dbRepositoryImageMap, nil
}
