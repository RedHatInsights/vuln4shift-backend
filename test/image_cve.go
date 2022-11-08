package test

import (
	"app/base/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAllImageCves(t *testing.T, imageID int64) (imageCves []models.ImageCve) {
	assert.Nil(t, DB.Where("image_id = ?", imageID).Order("image_id").Find(&imageCves).Error)
	return imageCves
}

func DeleteImageCves(t *testing.T, imageID int64) {
	assert.Nil(t, DB.Where("image_id = ?", imageID).Delete(&models.ImageCve{}).Error)
}
