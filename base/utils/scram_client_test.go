package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"testing"

	"github.com/xdg-go/scram"

	"github.com/stretchr/testify/assert"
)

func TestScramInitiate(t *testing.T) {
	for _, hf := range []scram.HashGeneratorFcn{sha256.New, sha512.New} {
		x := XDGSCRAMClient{HashGeneratorFcn: hf}
		assert.Nil(t, x.Begin("usr", "pswd", "authz-id"))
		assert.NotNil(t, x.Client)
		assert.NotNil(t, x.ClientConversation)

		// Initiate authentication conversation
		resp, err := x.Step("")
		assert.Nil(t, err)
		assert.NotEqual(t, "", resp)

		// Conversation should not complete yet
		assert.False(t, x.Done())
	}
}

func TestNewScramClientError(t *testing.T) {
	x := XDGSCRAMClient{HashGeneratorFcn: sha256.New}
	forbiddenRune := string(rune(568))
	assert.NotNil(t, x.Begin(forbiddenRune, "", ""))
}
