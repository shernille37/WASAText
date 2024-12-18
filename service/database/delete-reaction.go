package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) DeleteReaction(reactionID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) error {

	var res sql.Result

	const queryDeleteReaction = `
		DELETE FROM Reaction WHERE
		reactionID = ? AND messageID = ? AND 
		messageID IN (
			SELECT m.messageID
			FROM Message m
			WHERE m.conversationID = ? 
		);
	`

	res, err := db.c.Exec(queryDeleteReaction, reactionID, messageID, conversationID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Reaction does not exist")
	}

	return nil

}
