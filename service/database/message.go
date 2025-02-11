package database

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type MessageBody struct {
	ReplyMessageID *uuid.UUID
	Message        string
	Image          *string
}

type ForwardMessageBody struct {
	Source      uuid.UUID
	Destination *uuid.UUID
	ReceiverID  *uuid.UUID
}

func (db *appdbimpl) ListMessages(conversationID uuid.UUID) ([]Message, error) {
	var res []Message

	const queryMessages = `
		SELECT m.messageID, m.senderID, u.username,m.conversationID, m.timestamp, m.messageType,m.messageStatus, m.message, m.hasImage ,m.image,
		m.replyMessageID,u1.username, m1.message, u2.username
		FROM (
		( 
		(
		(Message m LEFT JOIN Message m1 ON m.replyMessageID = m1.messageID) 
		LEFT JOIN User u ON m.senderID = u.userID
		)
		LEFT JOIN User u1 ON m1.senderID = u1.userID )
		LEFT JOIN Message m2 ON m.forwardSourceMessageID = m2.messageID )
		LEFT JOIN User u2 ON m2.senderID = u2.userID

		WHERE m.conversationID = ?;
	`

	const queryReactions = `
		SELECT r.emoji, COUNT(r.emoji) as count
		FROM Reaction r
		WHERE r.messageID = ?
		GROUP BY r.emoji;
	`

	rows, err := db.c.Query(queryMessages, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m Message
		if err = rows.Scan(&m.MessageID, &m.SenderID, &m.SenderName, &m.ConversationID, &m.Timestamp, &m.MessageType, &m.MessageStatus, &m.Message, &m.HasImage, &m.Image, &m.ReplyMessageID, &m.ReplyRecipientName, &m.ReplyMessage, &m.ForwardedFromName); err != nil {
			return nil, err
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
		SELECT m.messageID, m.senderID, u.username ,m.conversationID, m.timestamp, m.messageType,m.messageStatus, m.message, m.hasImage, m.image,
		m.replyMessageID, u1.username,m1.message
		FROM (
		(Message m LEFT JOIN Message m1 ON m.replyMessageID = m1.messageID) 
		LEFT JOIN User u ON m.senderID = u.userID)
		LEFT JOIN User u1 ON m1.senderID = u1.userID
		WHERE m.conversationID = ? AND m.messageID = ?
	`

	tx, err := db.c.Begin()

	if err != nil {
		return res, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("%w -- Rollback Failed: %w", err, rbErr)
			}
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

	if err = tx.QueryRow(queryResponse, conversationID, messageID).Scan(&res.MessageID, &res.SenderID, &res.SenderName, &res.ConversationID, &res.Timestamp, &res.MessageType, &res.MessageStatus, &res.Message, &res.HasImage, &res.Image, &res.ReplyMessageID, &res.ReplyRecipientName, &res.ReplyMessage); err != nil {
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

		if err = rows.Scan(&u.UserID, &u.Username, &u.Image, &reader.Timestamp); err != nil {
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

func (db *appdbimpl) ForwardMessage(senderID uuid.UUID, messageID uuid.UUID, fmb ForwardMessageBody) (Conversation, error) {

	var res Conversation

	const queryMessage = `
		SELECT m.message, m.hasImage, m.image
		FROM Message m
		WHERE m.messageID = ? AND m.conversationID = ?; 
	`

	const queryExistingConversation = `
		SELECT c.conversationID
		FROM Conversation c, Members m1, Members m2
		WHERE c.conversationType = 'private' AND c.conversationID = m1.conversationID AND 
		m1.userID = ? AND c.conversationID = m2.conversationID AND 
		m2.userID = ?;
	`

	const queryForward = `
		INSERT INTO Message(messageID, senderID, conversationID, forwardSourceMessageID, messageType, hasImage, message, image)
		VALUES (?,?,?,?,?,?,?,?);
	`

	const queryAddConversation = `
	INSERT INTO Conversation(conversationID, conversationType) VALUES
	(?,'private');
	`

	const queryAddMembers = `
	INSERT INTO Members(userID, conversationID) VALUES
	(?,?);

	`

	var sourceMessage Message

	tx, err := db.c.Begin()
	if err != nil {
		return res, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("%w -- Rollback Failed: %w", err, rbErr)
			}
		}
	}()

	if err = tx.QueryRow(queryMessage, messageID, fmb.Source).Scan(&sourceMessage.Message, &sourceMessage.HasImage, &sourceMessage.Image); err != nil {
		return res, err
	}

	newMessageID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	sourceMessage.MessageType = "forward"
	var newConversationID uuid.UUID

	var conversationExist error
	if fmb.ReceiverID != nil {

		conversationExist = db.CheckExistingConversation(senderID, *fmb.ReceiverID)

		// Check if NO conversation already exists
		if conversationExist == nil {
			// Add new private conversation
			newConversationID, err = uuid.NewV4()

			if err != nil {
				return res, err
			}
			if _, err = tx.Exec(queryAddConversation, newConversationID); err != nil {
				return res, err
			}

			// Add Members
			mems := [2]uuid.UUID{senderID, *fmb.ReceiverID}
			for _, id := range mems {

				if _, err = tx.Exec(queryAddMembers, id, newConversationID); err != nil {
					return res, err
				}

			}

		}

	}

	var conversationID uuid.UUID
	// If a new conversation is Added
	if fmb.ReceiverID != nil && conversationExist == nil {
		conversationID = newConversationID
	} else if fmb.ReceiverID != nil && conversationExist != nil {
		// Query the existing private conversation
		if err = tx.QueryRow(queryExistingConversation, senderID, fmb.ReceiverID).Scan(&conversationID); err != nil {
			return res, err
		}
	} else {
		conversationID = *fmb.Destination
	}

	// Execute forward message (Add message)
	if _, err = tx.Exec(queryForward, newMessageID, senderID, conversationID, messageID, sourceMessage.MessageType, sourceMessage.HasImage, sourceMessage.Message, sourceMessage.Image); err != nil {
		return res, err
	}

	if err = tx.Commit(); err != nil {
		return res, err
	}

	res, err = db.GetConversation(senderID, conversationID)
	if err != nil {
		return res, err
	}
	res.Members = nil

	return res, nil

}
