package database

import (
	"errors"

	"github.com/gofrs/uuid"
)

type MessageBody struct {
	ReplyMessageID *uuid.UUID
	Message        string
	Image          *string
}

type ForwardMessageBody struct {
	Source      uuid.UUID `json:"source"`
	Destination uuid.UUID `json:"destination"`
}

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

func (db *appdbimpl) AddMessage(senderID uuid.UUID, conversationID uuid.UUID, mb MessageBody) (Message, error) {

	var res Message

	const queryAddMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, replyMessageID, messageType, hasImage, message, image)
		VALUES (?,?,?,?,?,?,?,?);
	`

	const queryResponse = `
		SELECT m.messageID, m.senderID, m.conversationID, m.timestamp, m.messageType,m.messageStatus, m.message, m.hasImage, m.image,
		m.replyMessageID, m1.message
		FROM Message m LEFT JOIN Message m1 ON m.replyMessageID = m1.messageID 
		WHERE m.conversationID = ? AND m.messageID = ?
	`

	tx, err := db.c.Begin()

	if err != nil {
		return res, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

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
	}

	if _, err = tx.Exec(queryAddMessage, messageID, senderID, conversationID, mb.ReplyMessageID, messageType, hasImage, mb.Message, mb.Image); err != nil {
		return res, err
	}

	if err = tx.QueryRow(queryResponse, conversationID, messageID).Scan(&res.MessageID, &res.SenderID, &res.ConversationID, &res.Timestamp, &res.MessageType, &res.MessageStatus, &res.Message, &res.HasImage, &res.Image, &res.ReplyMessageID, &res.ReplyMessage); err != nil {
		return res, err
	}


	if err = tx.Commit(); err != nil {
		return res, err
	}

	return res, nil

}

func (db *appdbimpl) DeleteMessage(conversationID uuid.UUID, messageID uuid.UUID) error {

	const queryDelete = `
		DELETE FROM Message
		WHERE conversationID = ? AND messageID = ?;
	`

	res, err := db.c.Exec(queryDelete, conversationID, messageID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Message does not exist")
	}

	return nil
}

func (db *appdbimpl) ListReaders(conversationID uuid.UUID, messageID uuid.UUID) ([]Reader, error) {

	var res []Reader

	const queryReaders = `
		SELECT u.userID, u.username, u.image, m.timestamp
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil

}

func (db *appdbimpl) ForwardMessage(senderID uuid.UUID, messageID uuid.UUID, fmb ForwardMessageBody) error {

	const queryMessage = `
		SELECT m.message, m.hasImage, m.image
		FROM Message m
		WHERE m.messageID = ? AND m.conversationID = ?; 
	`

	const queryForward = `
		INSERT INTO Message(messageID, senderID, conversationID, messageType, hasImage, message, image)
		VALUES (?,?,?,?,?,?,?);
	`

	var sourceMessage Message

	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()



	if err = tx.QueryRow(queryMessage, messageID, fmb.Source).Scan(&sourceMessage.Message, &sourceMessage.HasImage, &sourceMessage.Image); err != nil {
		return err
	}

	newMessageID, err := uuid.NewV4()

	if err != nil {
		return err
	}

	sourceMessage.MessageType = "forward"

	if !sourceMessage.HasImage {
		sourceMessage.Image = nil
	}

	if _, err = tx.Exec(queryForward, newMessageID, senderID, fmb.Destination, sourceMessage.MessageType, sourceMessage.HasImage, sourceMessage.Message, sourceMessage.Image); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}


