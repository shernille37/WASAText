package database

import (
	"errors"

	"github.com/gofrs/uuid"
)

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
		return errors.New("Conversation does not exist")
	}

	return nil

}
