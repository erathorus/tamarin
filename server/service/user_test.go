package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lattetalk/lattetalk/db"
	"gitlab.com/lattetalk/lattetalk/model"
)

func testCreateUser(t *testing.T) {
	user := &model.User{
		AuthID:         "Auth0 id",
		FamilyName:     "Family",
		GivenName:      "Given",
		ProfilePicture: "Picture",
	}
	id := CreateUser(user)

	user2 := GetUserByID(id)

	assert.Equal(t, id, user2.ID)
	assert.Equal(t, user.AuthID, user2.AuthID)
	assert.Equal(t, user.FamilyName, user2.FamilyName)
	assert.Equal(t, user.GivenName, user2.GivenName)
	assert.Equal(t, user.ProfilePicture, user2.ProfilePicture)
}

func TestUserService(t *testing.T) {
	t.Run("CreateUser", userWrapper(testCreateUser))
}

func userWrapper(f func(*testing.T)) func(*testing.T) {
	truncateUserTable()
	return f
}

func truncateUserTable() {
	if _, err := db.DB.Exec("TRUNCATE application_user CASCADE"); err != nil {
		panic(err)
	}
}
