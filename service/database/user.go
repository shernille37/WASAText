package database

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/mattn/go-sqlite3"
	"github.com/shernille37/WASAText/service/api/constants"
)

func (db *appdbimpl) ListUsers(id uuid.UUID) ([]User, error) {

	var res []User

	const queryUsers = `
		SELECT u.userID, u.username, u.image
		FROM User u
		WHERE u.userID <> ?;
	`

	rows, err := db.c.Query(queryUsers, id)

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

func (db *appdbimpl) UpdateUsername(userID uuid.UUID, newUsername string) error {

	const queryUpdate = `
		UPDATE User
		SET username = ?
		WHERE userID = ?;
	`

	res, err := db.c.Exec(queryUpdate, newUsername, userID)

	if err != nil {
		// Check if the error is due to a UNIQUE constraint violation
		// Username already exists
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint && strings.HasPrefix(err.Error(), "UNIQUE") {
			return errors.New(constants.USER_ALREADY_EXISTS)
		}

		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New(constants.NO_USER)
	}

	return nil

}


func (db *appdbimpl) UpdateUserImage(userID uuid.UUID, newUserImage string) error {

	const queryUpdate = `
		UPDATE User 
		SET image = ?
		WHERE userID = ?;
	`

	res, err := db.c.Exec(queryUpdate, newUserImage, userID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New(constants.NO_USER)
	}

	return nil

}