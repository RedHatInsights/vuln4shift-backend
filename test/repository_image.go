package test

import (
	"app/base/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRepoImage(t *testing.T, image models.RepositoryImage) {
	assert.Nil(t, DB.Create(image).Error)
}

func GetRepoImages(t *testing.T, repoID int64) (repoImages []models.RepositoryImage) {
	assert.Nil(t, DB.Where("repository_id = ?", repoID).Find(&repoImages).Error)
	return repoImages
}

func DeleteRepoImageByImageID(t *testing.T, imageID int64) {
	assert.Nil(t, DB.Where("image_id = ?", imageID).Delete(&models.RepositoryImage{}).Error)
}

func GetClusterRepoImages(t *testing.T, id int64) (images []models.RepositoryImage) {
	assert.Nil(t, DB.Model(models.RepositoryImage{}).
		Joins("JOIN cluster_image ON repository_image.image_id = cluster_image.image_id").
		Joins("JOIN image_cve ON cluster_image.image_id = image_cve.image_id").
		Joins("JOIN cluster ON cluster_image.cluster_id = cluster.id").
		Order("repository_image.repository_id").
		Where("cluster.id = ?", id).
		Scan(&images).Error)
	return images
}
