package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) DeleteMessage(conversationID uuid.UUID, messageID uuid.UUID, typeFlag bool, userID uuid.UUID) error {

	const queryDeleteAll = `
		DELETE FROM Message
		WHERE conversationID = ? AND messageID = ?;
	`

	const queryDeleteOne = `
		UPDATE UserMessageStatus 
		SET visibility = 0
		WHERE userID = ? AND messageID = ?;
	`

	const queryUserMessageStatus = `
		SELECT ums.visibility
		FROM UserMessageStatus ums
		WHERE ums.messageID = ? AND userID <> ?;
	`

	var res sql.Result
	var err error

	// Delete ALL
	if typeFlag {
		if res, err = db.c.Exec(queryDeleteAll, conversationID, messageID); err != nil {
			return err
		}

	} else {

		// If all of the visibilites of the message is FALSE then delete all the message
		var visibilities []bool

		rows, err := db.c.Query(queryUserMessageStatus, messageID, userID)

		if err != nil {
			return nil
		}

		defer rows.Close()

		for rows.Next() {
			var vis bool
			if err = rows.Scan(&vis); err != nil {
				return nil
			}

			visibilities = append(visibilities, vis)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		deleteAll := true
		for _, vis := range visibilities {
			if vis {
				deleteAll = false
			}
		}

		if deleteAll {
			// Delete ALL
			if res, err = db.c.Exec(queryDeleteAll, conversationID, messageID); err != nil {
				return err
			}
		} else {
			// Delete ONE
			if res, err = db.c.Exec(queryDeleteOne, userID, messageID); err != nil {
				return err
			}
		}

	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Message does not exist")
	}

	return nil
}
