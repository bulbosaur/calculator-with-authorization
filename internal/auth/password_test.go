package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHashAndCompare(t *testing.T) {
	password := "securePassword123"

	hash, err := testingAuthService().GenerateHash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	assert.True(t, testingAuthService().Compare(hash, password))

	assert.False(t, testingAuthService().Compare(hash, "wrongPassword"))
}

func TestGenerateHash_EmptyPassword(t *testing.T) {
	hash, err := testingAuthService().GenerateHash("")
	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestCompare_InvalidHash(t *testing.T) {
	testCases := []struct {
		name     string
		hash     string
		password string
		expected bool
	}{
		{
			name:     "empty hash",
			hash:     "",
			password: "password",
			expected: false,
		},
		{
			name:     "malformed hash",
			hash:     "invalid_hash",
			password: "password",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, testingAuthService().Compare(tc.hash, tc.password))
		})
	}
}

func TestCompare_EmptyPassword(t *testing.T) {
	hash, err := testingAuthService().GenerateHash("password")
	assert.NoError(t, err)

	assert.False(t, testingAuthService().Compare(hash, ""))
}
