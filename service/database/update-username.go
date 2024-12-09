package database

import (
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) UpdateUsername(userID uuid.UUID, newUsername string) error {

	const queryUpdate = `
		UPDATE User
		SET username = ?
		WHERE userID = ?;
	`

	res, err := db.c.Exec(queryUpdate, newUsername, userID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Conversation does not exist")
	}

	return nil

}
