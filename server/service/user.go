package service

import (
	"database/sql"

	"gitlab.com/lattetalk/lattetalk/db"
	"gitlab.com/lattetalk/lattetalk/model"
)

// CreateUser creates a new user. If the user is already exist
// the id return is set to -1.
func CreateUser(user *model.User) int64 {
	stmt, err := db.DB.Prepare(
		`INSERT INTO application_user (auth_id, family_name, given_name, profile_picture)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING
				RETURNING id`,
	)
	if err != nil {
		panic(err)
	}
	var id int64
	err = stmt.QueryRow(user.AuthID, user.FamilyName, user.GivenName, user.ProfilePicture).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
		panic(err)
	}
	return id
}

func GetUserByOtherUserID(id, otherID int64) *model.User {
	if CheckUsersAreFriendsByID(id, otherID) {
		return GetUserByID(id)
	}
	return GetUserPublicProfile(id)
}

func GetUserPublicProfile(id int64) *model.User {
	stmt, err := db.DB.Prepare(
		`SELECT family_name, given_name
				FROM application_user
				WHERE id=$1`,
	)
	if err != nil {
		panic(err)
	}

	user := new(model.User)
	user.ID = id
	if err = stmt.QueryRow(id).Scan(&user.FamilyName, &user.GivenName); err != nil {
		panic(err)
	}
	return user
}

func CheckUsersAreFriendsByID(id1, id2 int64) bool {
	if id1 == id2 {
		return true
	}

	if id1 > id2 {
		tmp := id1
		id1 = id2
		id2 = tmp
	}

	stmt, err := db.DB.Prepare(
		`SELECT EXISTS (SELECT 1 FROM friendship WHERE user1_id=$1 AND user2_id=$2)`,
	)
	if err != nil {
		panic(err)
	}

	var exist bool
	if err = stmt.QueryRow(id1, id2).Scan(&exist); err != nil {
		panic(err)
	}

	return exist
}

func GetUserByAuthID(authID string) *model.User {
	stmt, err := db.DB.Prepare(
		`SELECT id, family_name, given_name, profile_picture
				FROM application_user
				WHERE auth_id=$1`,
	)
	if err != nil {
		panic(err)
	}

	user := new(model.User)
	user.AuthID = authID
	if err = stmt.QueryRow(authID).Scan(&user.ID, &user.FamilyName, &user.GivenName, &user.ProfilePicture); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return user
}

func GetUserIDByAuthID(authID string) int64 {
	stmt, err := db.DB.Prepare(
		`SELECT id
				FROM application_user
				WHERE auth_id=$1`,
	)
	if err != nil {
		panic(err)
	}

	var id int64
	if err = stmt.QueryRow(authID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
		panic(err)
	}

	return id
}

func GetUserByID(id int64) *model.User {
	stmt, err := db.DB.Prepare(
		`SELECT family_name, given_name, profile_picture
				FROM application_user
				WHERE id=$1`,
	)
	if err != nil {
		panic(err)
	}

	user := new(model.User)
	user.ID = id
	if err = stmt.QueryRow(id).Scan(&user.FamilyName, &user.GivenName,
		&user.ProfilePicture); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return user
}

func CheckUserExistByAuthID(authID string) bool {
	stmt, err := db.DB.Prepare(
		`SELECT EXISTS (SELECT 1 FROM application_user WHERE auth_id=$1)`,
	)
	if err != nil {
		panic(err)
	}
	var exist bool
	if err = stmt.QueryRow(authID).Scan(&exist); err != nil {
		panic(err)
	}
	return exist
}
