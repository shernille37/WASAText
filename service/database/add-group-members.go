package database

import (
	"github.com/gofrs/uuid"
)

type GroupMemberBody struct {
	Members []uuid.UUID
}

func (db *appdbimpl) AddGroupMembers(conversationID uuid.UUID, gmb GroupMemberBody) ([]User, error) {

	var res []User

	const queryAddMembers = `
		INSERT INTO Members(userID, conversationID)
		VALUES (?, ?);
	`

	const queryResponse = `
	SELECT u.userID, u.username, u.image 
	FROM User u
	WHERE u.userID = ?;
	`

	for _, member := range gmb.Members {

		_, err := db.c.Exec(queryAddMembers, member, conversationID)

		if err != nil {
			return nil, err
		}
	}

	for _, member := range gmb.Members {
		var u User
		if err := db.c.QueryRow(queryResponse, member).Scan(&u.UserID, &u.Name, &u.Image); err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	return res, nil

}
