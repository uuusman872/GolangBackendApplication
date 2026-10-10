package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser(t *testing.T) {
	defer cleanupTestData(t)
	repo := NewSQLUserRepository(testDB)
	userID, err := repo.CreateUser("john Doe", "john@doe.com", "password", "avatar.jpg")
	assert.NoError(t, err)
	assert.Greater(t, userID, 0)

	user, err := repo.GetUserByEmail("john@doe.com")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)

	userID, err = repo.CreateUser(
		"job Doe",
		"john2@doe.com",
		generateString(73),
		"avatar",
	)
	assert.Error(t, err)
	assert.Equal(t, userID, 0)
}

func TestSQLUserRepository_CreateUser_DuplicateEmail(t *testing.T) {
	defer cleanupTestData(t)
	repo := NewSQLUserRepository(testDB)
	userID, err := repo.CreateUser(
		"john Doe",
		"john@doe.com",
		"testpassword",
		"avatar",
	)
	assert.Nil(t, err)
	assert.Greater(t, userID, 0)

	userID, err = repo.CreateUser(
		"john Doe",
		"john@doe.com",
		"testpassword",
		"avatar",
	)
	assert.Error(t, err)
}

func TestSQLUserRepository_Authenticate_Success(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)
	userID, err := repo.CreateUser("john Doe", "john@doe.com", "testpassword", "avatar")
	assert.Nil(t, err)
	assert.Greater(t, userID, 0)

	userID, err = repo.Authenticate("john@doe.com", "testpassword")
	assert.NoError(t, err)
	_, err = repo.Authenticate("john@doe.com", "wrongpassword")
	assert.Error(t, err)
}

func generateString(n int) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = 'a'
	}
	return string(buf)
}
