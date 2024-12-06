package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListMessages(conversationID uuid.UUID) ([]Message, error) {
	var res []Message

	const queryMessages = `
		SELECT m.messageID, m.senderID, m.conversationID, m.timestamp, m.messageType, m.messageStatus, m.timeDelivered, COALESCE(m.message, ''), COALESCE(m.image, '')
		FROM Message m
		WHERE m.conversationID = ?
	`

	rows, err := db.c.Query(queryMessages, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m Message

		if err = rows.Scan(&m.MessageID, &m.SenderID, &m.ConversationID, &m.Timestamp, &m.MessageType, &m.MessageStatus, &m.TimeDelivered, &m.Message, &m.Image); err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	return res, nil
}
