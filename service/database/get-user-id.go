package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserByID(id string) (User, error) {

	var res User
	const queryUser = `
		SELECT u.userID, u.username, u.image
		FROM User u
		WHERE u.userID = ?;
	`

	err := db.c.QueryRow(queryUser, id).Scan(&res.UserID, &res.Name, &res.Image)

	// No user found
	if errors.Is(err, sql.ErrNoRows) {
		return res, err
	}

	return res, nil

}
