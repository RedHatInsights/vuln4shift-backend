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
