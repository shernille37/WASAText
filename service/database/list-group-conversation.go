package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListGroupConversation(id uuid.UUID) ([]GroupConversation, error) {

	var res []GroupConversation

	const queryGroupConversation = `SELECT c.conversationID, c.groupName, c.groupImage 
	FROM Conversation c, Members m, User u
	WHERE c.conversationType = 'group' AND u.userID = ? AND c.conversationID = m.conversationID AND m.userID = ?;
	`
	const queryLatestMessage = `
		SELECT	m.messageType, m.timestamp, COALESCE(m.message, '') FROM Message m 
		WHERE m.conversationID = ?
		ORDER BY m.timestamp
		LIMIT 1;
	`

	rows, err := db.c.Query(queryGroupConversation, id, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var gc GroupConversation
		var lm LatestMessage

		err = rows.Scan(&gc.ConversationID, &gc.GroupName, &gc.GroupImage)
		if err != nil {
			return nil, err
		}

		// Fetch the latest message
		err = db.c.QueryRow(queryLatestMessage, gc.ConversationID).Scan(&lm.MessageType, &lm.Timestamp, &lm.Message)
		
		if err != nil && err != sql.ErrNoRows {
			
			return nil, err
		}
		gc.LatestMessage = &lm

		res = append(res, gc)
	}

	return res, nil

}