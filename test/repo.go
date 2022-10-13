package test

import (
	"app/base/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRepo(t *testing.T, repo models.Repository) {
	result := DB.Create(&repo)
	assert.Nil(t, result.Error)
	assert.True(t, result.RowsAffected > 0)
}

func DeleteRepo(t *testing.T, registry, repository string) {
	assert.Nil(t, DB.Where("repository = ? AND registry = ?", repository, registry).Delete(&models.Repository{}).Error)
}

func GetRepo(t *testing.T, repository string) (repo models.Repository) {
	result := DB.Model(models.Repository{}).Where("repository = ?", repository).Find(&repo)
	assert.Nil(t, result.Error)
	assert.True(t, result.RowsAffected > 0)
	return repo
}
