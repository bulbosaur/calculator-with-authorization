package auth_test

import (
	"testing"
	"time"

	"github.com/bulbosaur/calculator-with-authorization/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testingService() *auth.Service {
	return &auth.Service{
		SecretKey:     "test_secret_key",
		TokenDuration: 24 * time.Hour,
	}
}

func TestGenerateAndParseJWT(t *testing.T) {
	viper.Set("jwt.token_duration", 24)
	userID := 123

	token, err := testingService().GenerateJWT(userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := testingService().ParseJWT(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Minute)
}

func TestParseJWT_InvalidToken(t *testing.T) {
	expiredToken := func() string {
		claims := &auth.Claims{
			UserID: 123,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, _ := token.SignedString([]byte("test_secret_key"))
		return s
	}()

	validToken, err := testingService().GenerateJWT(123)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		token      string
		Service    *auth.Service
		expectErr  bool
		errMessage string
	}{
		{
			name:       "invalid signature",
			token:      validToken[:len(validToken)-5] + "XXXXX",
			Service:    testingService(),
			expectErr:  true,
			errMessage: "signature is invalid",
		},
		{
			name:       "expired token",
			token:      expiredToken,
			Service:    testingService(),
			expectErr:  true,
			errMessage: "token is expired",
		},
		{
			name:       "malformed token",
			token:      "malformed.token",
			Service:    testingService(),
			expectErr:  true,
			errMessage: "token contains an invalid number of segments",
		},
		{
			name:  "wrong secret key",
			token: validToken,
			Service: &auth.Service{
				SecretKey:     "wrong_secret_key",
				TokenDuration: 24 * time.Hour,
			},
			expectErr:  true,
			errMessage: "signature is invalid",
		},
		{
			name:       "empty token",
			token:      "",
			Service:    testingService(),
			expectErr:  true,
			errMessage: "token contains an invalid number of segments",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := tc.Service.ParseJWT(tc.token)
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
	Service := &auth.Service{
		SecretKey:     "",
		TokenDuration: 24 * time.Hour,
	}
	token, err := Service.GenerateJWT(123)
	require.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "secret key is empty")
}
