package database

import (
	"github.com/gofrs/uuid"
)

type MessageBody struct {
	ReplyMessageID *uuid.UUID
	Message        string
	Image          *string
}

func (db *appdbimpl) AddMessage(senderID uuid.UUID, conversationID uuid.UUID, mb MessageBody) (Message, error) {

	var res Message

	const queryAddMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, replyMessageID, messageType, hasImage, message, image)
		VALUES (?,?,?,?,?,?,?,?);
	`

	const queryResponse = `
		SELECT m.messageID, m.senderID, m.conversationID, m.timestamp, m.messageType,m.messageStatus, m.message, m.hasImage,COALESCE(m.image, ''),
		m.replyMessageID ,COALESCE(m1.message, '')
		FROM Message m LEFT JOIN Message m1 ON m.replyMessageID = m1.messageID 
		WHERE m.conversationID = ? AND m.messageID = ?
	`

	// ADD Message
	messageID, err := uuid.NewV4()

	if err != nil {
		return res, err
	}

	var messageType string
	var hasImage bool

	if mb.ReplyMessageID != nil {
		messageType = "reply"
	} else {
		messageType = "default"
	}

	if mb.Image != nil {
		hasImage = true
	} else {
		hasImage = false
	}

	_, err = db.c.Exec(queryAddMessage, messageID, senderID, conversationID, mb.ReplyMessageID, messageType, hasImage, mb.Message, mb.Image)

	if err != nil {
		return res, nil
	}

	if err := db.c.QueryRow(queryResponse, conversationID, messageID).Scan(&res.MessageID, &res.SenderID, &res.ConversationID, &res.Timestamp, &res.MessageType, &res.MessageStatus, &res.Message, &res.HasImage, &res.Image, &res.ReplyMessageID, &res.ReplyMessage); err != nil {
		return res, nil
	}

	if !res.HasImage {
		res.Image = nil
	}

	if res.ReplyMessageID == nil {
		res.ReplyMessage = nil
	}

	return res, nil

}
