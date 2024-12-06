package database

import (
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) UpdateGroupImage(conversationID uuid.UUID, newGroupPhoto string) error {

	const queryUpdate = `
		UPDATE Conversation
		SET groupImage = ?
		WHERE conversationID = ? AND conversationType = 'group';
	`

	res, err := db.c.Exec(queryUpdate, newGroupPhoto, conversationID)

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
