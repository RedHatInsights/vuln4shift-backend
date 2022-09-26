package test

import (
	"app/base/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetAccounts(t *testing.T) (accounts []models.Account) {
	result := DB.Model(models.Account{}).Scan(&accounts)
	assert.Nil(t, result.Error)
	assert.True(t, len(accounts) > 0)
	return accounts
}
