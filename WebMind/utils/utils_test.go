package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	username := "testuser"
	// 生成 JWT
	token, err := GenerateToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 验证 JWT
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, username, claims.Username)
}
