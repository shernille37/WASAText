package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListPrivateConversation(id uuid.UUID) ([]PrivateConversation, error) {
	var res []PrivateConversation

	const queryPrivateConversation = `
		SELECT c.conversationID
		FROM Conversation c, Members m
		WHERE c.conversationType = 'personal' AND c.conversationID = m.conversationID AND m.userID = ?;
	`

	const queryChatmate = `
		SELECT u.userID, u.username, u.image
		FROM User u, Members m 
		WHERE m.conversationID = ? AND m.userID = u.userID AND m.userID <> ?;
	`

	const queryLatestMessage = `
		SELECT m.messageType, m.timestamp, COALESCE(m.message, '') 
		FROM Message m 
		WHERE m.conversationID = ?
		ORDER BY m.timestamp DESC
		LIMIT 1;
	`

	rows, err := db.c.Query(queryPrivateConversation, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		var pc PrivateConversation
		var lm LatestMessage

		if err = rows.Scan(&pc.ConversationID); err != nil {
			return nil, err
		}

		// Fetch chatmate
		if err = db.c.QueryRow(queryChatmate, pc.ConversationID, id).Scan(&u.UserID, &u.Name, &u.Image); err != nil {
			return nil, err
		}

		// Fetch the latest message
		if err = db.c.QueryRow(queryLatestMessage, pc.ConversationID).Scan(&lm.MessageType, &lm.Timestamp, &lm.Message); err != nil {
			return nil, err
		}

		pc.User = &u
		pc.LatestMessage = &lm
		res = append(res, pc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
