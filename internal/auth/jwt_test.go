package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseJWT(t *testing.T) {
	viper.Set("jwt.token_duration", 24)
	userID := 123
	secretKey := "test_secret_key"

	token, err := GenerateJWT(userID, secretKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseJWT(token, secretKey)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Minute)
}

func TestParseJWT_InvalidToken(t *testing.T) {
	expiredToken := func() string {
		claims := &Claims{
			UserID: 123,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, _ := token.SignedString([]byte("test_secret_key"))
		return s
	}()

	validToken, err := GenerateJWT(123, "test_secret_key")
	require.NoError(t, err)

	testCases := []struct {
		name       string
		token      string
		secretKey  string
		expectErr  bool
		errMessage string
	}{
		{
			name:       "invalid signature",
			token:      validToken[:len(validToken)-5] + "XXXXX",
			secretKey:  "test_secret_key",
			expectErr:  true,
			errMessage: "signature is invalid",
		},
		{
			name:       "expired token",
			token:      expiredToken,
			secretKey:  "test_secret_key",
			expectErr:  true,
			errMessage: "token is expired",
		},
		{
			name:       "malformed token",
			token:      "malformed.token",
			secretKey:  "test_secret_key",
			expectErr:  true,
			errMessage: "token contains an invalid number of segments",
		},
		{
			name:       "wrong secret key",
			token:      validToken,
			secretKey:  "wrong_secret_key",
			expectErr:  true,
			errMessage: "signature is invalid",
		},
		{
			name:       "empty token",
			token:      "",
			secretKey:  "test_secret_key",
			expectErr:  true,
			errMessage: "token contains an invalid number of segments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := ParseJWT(tc.token, tc.secretKey)
			if tc.expectErr {
				require.Error(t, err)
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)
				}
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestGenerateJWT_EmptySecret(t *testing.T) {
	token, err := GenerateJWT(123, "")
	assert.Error(t, err)
	assert.Empty(t, token)
}
