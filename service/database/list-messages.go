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

	const queryReactions = `
		SELECT r.emoji, COUNT(r.emoji) as count
		FROM Reaction r, User u, Message m
		WHERE r.messageID = ? AND r.messageID = m.messageID AND 
		r.userID = u.userID AND m.conversationID = ?
		GROUP BY r.emoji;
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

		if !m.HasImage {
			m.Image = nil
		}

		if m.ReplyMessageID == nil {
			m.ReplyMessage = nil
		}

		rowReactions, err := db.c.Query(queryReactions, m.MessageID, m.ConversationID)

		if err != nil {
			return nil, err
		}

		defer rowReactions.Close()

		for rowReactions.Next() {
			var r Reaction
			r.ReactionID = nil
			if err = rowReactions.Scan(&r.Unicode, &r.Count); err != nil {
				return nil, err
			}

			m.Reactions = append(m.Reactions, r)
		}

		if err = rowReactions.Err(); err != nil {
			return nil, err
		}

		res = append(res, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
