package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) Login(user User) (User, error) {

	var res User
	const queryUser = `
		SELECT u.userID, u.username, u.image
		FROM User u
		WHERE u.username = ?;
	`

	const queryRegisterUser = `
		INSERT INTO User(userID, username, image) VALUES 
		(?,?,"");
	`

	err := db.c.QueryRow(queryUser, user.Name).Scan(&res.UserID, &res.Name, &res.Image)

	// No user found -- Register the User to the system
	if errors.Is(err, sql.ErrNoRows) {
		userID, err := uuid.NewV4()
		if err != nil {
			return res, err
		}

		if _, err = db.c.Exec(queryRegisterUser, userID.String(), user.Name); err != nil {
			return res, err
		}

		res.UserID = userID
		res.Name = user.Name
		res.Image = nil

		return res, nil
	}

	return res, nil

}
