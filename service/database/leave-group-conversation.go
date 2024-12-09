package database

import (
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) LeaveGroupConversation(userID uuid.UUID, conversationID uuid.UUID) error {

	const queryDelete = `
		DELETE FROM Members
		WHERE conversationID = ? AND userID = ?;
	`

	res, err := db.c.Exec(queryDelete, conversationID, userID)

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
