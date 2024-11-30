package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListPrivateConversation(id uuid.UUID) ([]PrivateConversation, error) {
	var res []PrivateConversation

	const queryPrivateConversation = `
		SELECT c.conversationID, u.userID, u.username, u.image
		FROM Conversation c, Members m, User u
		WHERE c.conversationType = 'personal' AND u.userID = ? AND c.conversationID = m.conversationID AND m.userID = ?;
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
		var pc PrivateConversation
		var lm LatestMessage

		err = rows.Scan(&pc.ConversationID, &pc.UserID, &pc.Name, &pc.Image)
		if err != nil {
			return nil, err
		}

		// Fetch the latest message
		err = db.c.QueryRow(queryLatestMessage, pc.ConversationID).Scan(&lm.MessageType, &lm.Timestamp, &lm.Message)
		if err != nil && err != sql.ErrNoRows {
			
			return nil, err
		}

		pc.LatestMessage = &lm 
		res = append(res, pc)  
	}

	return res, nil
}
