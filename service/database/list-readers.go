package database

import "github.com/gofrs/uuid"

func (db *appdbimpl) ListReaders(conversationID uuid.UUID, messageID uuid.UUID) ([]Reader, error) {

	var res []Reader

	const queryReaders = `
		SELECT u.userID, u.username, u.image, timestamp
		FROM User u, Reader r, Conversation c, Message m
		WHERE c.conversationID = ? AND m.conversationID = c.conversationID AND
		m.messageID = ? AND m.messageID = r.messageID AND r.userID = u.userID;
	`

	rows, err := db.c.Query(queryReaders, conversationID, messageID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var reader Reader
		var u User

		if err = rows.Scan(&u.UserID, &u.Name, &u.Image, &reader.Timestamp); err != nil {
			return nil, err
		}

		reader.User = &u

		res = append(res, reader)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return res, nil

}
