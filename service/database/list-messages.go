package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListMessages(conversationID uuid.UUID) ([]Message, error) {
	var res []Message

	const queryMessages = `
		SELECT m.messageID, m.senderID, m.conversationID, m.timestamp, m.messageType,m.messageStatus, m.message, m.hasImage ,COALESCE(m.image, ''),
		m.replyMessageID ,COALESCE(m1.message, '')
		FROM Message m LEFT JOIN Message m1 ON m.replyMessageID = m1.messageID
		WHERE m.conversationID = ?
	`

	rows, err := db.c.Query(queryMessages, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m Message
		if err = rows.Scan(&m.MessageID, &m.SenderID, &m.ConversationID, &m.Timestamp, &m.MessageType, &m.MessageStatus, &m.Message, &m.HasImage, &m.Image, &m.ReplyMessageID, &m.ReplyMessage); err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
