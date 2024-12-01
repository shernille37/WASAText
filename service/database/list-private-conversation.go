package database

import (
	"database/sql"

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

	rows, err := db.c.Query(queryPrivateConversation, id, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		var pc PrivateConversation
		var lm LatestMessage

		err = rows.Scan(&pc.ConversationID)
		if err != nil {
			return nil, err
		}

		// Fetch chatmate
		err = db.c.QueryRow(queryChatmate, pc.ConversationID, id).Scan(&u.UserID, &u.Name, &u.Image)

		if err != nil {
			return nil ,err
		}

		// Fetch the latest message
		err = db.c.QueryRow(queryLatestMessage, pc.ConversationID).Scan(&lm.MessageType, &lm.Timestamp, &lm.Message)
		if err != nil && err != sql.ErrNoRows {
			
			return nil, err
		}

		pc.User = &u
		pc.LatestMessage = &lm 
		res = append(res, pc)  
	}

	return res, nil
}
