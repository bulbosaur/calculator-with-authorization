package repository_test

import (
	"testing"

	"github.com/bulbosaur/calculator-with-authorization/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetUser(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	user := &models.User{
		Login:        "testuser",
		PasswordHash: "hashedpassword",
	}

	userID, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.Greater(t, userID, 0)

	dbUser, err := repo.GetUserByLogin("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Login, dbUser.Login)
	assert.Equal(t, user.PasswordHash, dbUser.PasswordHash)
}

func TestGetNonexistentUser(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	_, err := repo.GetUserByLogin("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, models.ErrorUserNotFound, err)
}

func TestCreateUser_DuplicateLogin(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	user1 := &models.User{
		Login:        "testuser",
		PasswordHash: "hash1",
	}
	_, err := repo.CreateUser(user1)
	assert.NoError(t, err)

	user2 := &models.User{
		Login:        "testuser",
		PasswordHash: "hash2",
	}
	_, err = repo.CreateUser(user2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}
