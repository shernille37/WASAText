package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListGroupMembers(conversationID uuid.UUID) ([]User, error) {
	var res []User

	const queryMembers = `
		SELECT u.userID, u.username, u.image FROM 
		Members m, User u
		WHERE m.conversationID = ? AND m.userID = u.userID;
	`

	rows, err := db.c.Query(queryMembers, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		if err = rows.Scan(&u.UserID, &u.Name, &u.Image); err != nil {
			return nil, err
		}

		res = append(res, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
