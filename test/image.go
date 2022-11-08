package test

import (
	"app/base/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAllImages(t *testing.T) (images []models.Image) {
	assert.Nil(t, DB.Model(models.Image{}).Find(&images).Error)
	return images
}

func GetImageByID(t *testing.T, id int64) (img models.Image) {
	assert.Nil(t, DB.Model(models.Image{}).Where("id = ?", id).Find(&img).Error)
	return img
}

func GetImageByPyxisID(t *testing.T, id string) (img models.Image) {
	assert.Nil(t, DB.Model(models.Image{}).Where("pyxis_id = ?", id).Find(&img).Error)
	return img
}

func DeleteImage(t *testing.T, pyxisID string) {
	assert.Nil(t, DB.Where("pyxis_id = ?", pyxisID).Delete(&models.Image{}).Error)
}

func CreateImage(t *testing.T, image models.Image) {
	assert.Nil(t, DB.Create(&image).Error)
}
