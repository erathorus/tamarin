package service

import (
	"time"

	"gitlab.com/lattetalk/lattetalk/db"
	"gitlab.com/lattetalk/lattetalk/model"
)

func GetUserFriendsByID(id int64) []model.User {
	stmt, err := db.DB.Prepare(
		`SELECT id, family_name, given_name, profile_picture
				FROM application_user
				WHERE id IN (
					SELECT user1_id FROM friendship WHERE user2_id = $1
					UNION
					SELECT user2_id FROM friendship WHERE user1_id = $1
				)`,
	)
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}

	friends := make([]model.User, 0)
	var user model.User
	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.FamilyName, &user.GivenName, &user.ProfilePicture); err != nil {
			panic(err)
		}
		friends = append(friends, user)
	}

	return friends
}

func GetUserNewFriendsWithSimilarName(id int64, name string) []model.User {
	stmt, err := db.DB.Prepare(
		`SELECT id, family_name, given_name, profile_picture
				FROM application_user A
				WHERE 
					(given_name ILIKE $2||'%' 
					OR family_name ILIKE $2||'%' 
					OR given_name||' '||family_name ILIKE $2||'%')
					AND NOT EXISTS (
						SELECT 1 
						FROM friendship 
						WHERE 
							(user1_id=A.id AND user2_id=$1)
							OR (user2_id=A.id AND user1_id=$1)
					)
					AND id!=$1				
				LIMIT 10`,
	)
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Query(id, name)
	if err != nil {
		panic(err)
	}

	friends := make([]model.User, 0)
	var user model.User
	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.FamilyName, &user.GivenName, &user.ProfilePicture); err != nil {
			panic(err)
		}
		friends = append(friends, user)
	}

	return friends
}

func ShouldAddNewFriend(user1ID, user2ID int64) (*model.Conversation, error) {
	if user1ID > user2ID {
		tmp := user1ID
		user1ID = user2ID
		user2ID = tmp
	}
	now := time.Now()
	stmt, err := db.DB.Prepare(
		`INSERT INTO friendship (user1_id, user2_id, created_at) VALUES ($1, $2, $3)`,
	)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(user1ID, user2ID, now)
	if err != nil {
		return nil, err
	}
	stmt, err = db.DB.Prepare(
		`INSERT INTO conversation (created_at, updated_at) VALUES ($1, $2) RETURNING id`,
	)
	if err != nil {
		return nil, err
	}
	conversation := new(model.Conversation)
	if err = stmt.QueryRow(now, now).Scan(&conversation.ID); err != nil {
		return nil, err
	}
	stmt, err = db.DB.Prepare(
		`INSERT INTO user_conversation (user_id, conversation_id) VALUES ($1, $2)`,
	)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(user1ID, conversation.ID)
	if err != nil {
		return nil, err
	}
	stmt, err = db.DB.Prepare(
		`INSERT INTO user_conversation (user_id, conversation_id) VALUES ($1, $2)`,
	)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(user2ID, conversation.ID)
	if err != nil {
		return nil, err
	}
	conversation.UserIDs = append(conversation.UserIDs, user1ID)
	conversation.UserIDs = append(conversation.UserIDs, user2ID)
	conversation.CreatedAt = now
	conversation.UpdatedAt = now
	return conversation, nil
}
