package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

type MessagePrivateBody struct {
	ReceiverID uuid.UUID
	Message    string
	Image      *string
}

func (db *appdbimpl) ListPrivateConversation(id uuid.UUID) ([]Conversation, error) {
	var res []Conversation

	const queryPrivateConversation = `
		SELECT c.conversationID, c.conversationType
		FROM Conversation c, Members m
		WHERE c.conversationType = 'private' AND c.conversationID = m.conversationID AND m.userID = ?;
	`

	const queryChatmate = `
		SELECT u.userID, u.username, u.image
		FROM User u, Members m 
		WHERE m.conversationID = ? AND m.userID = u.userID AND m.userID <> ?;
	`

	rows, err := db.c.Query(queryPrivateConversation, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Conversation
		var pc PrivateConversation
		var u User

		if err = rows.Scan(&c.ConversationID, &c.Type); err != nil {
			return nil, err
		}

		// Fetch chatmate
		if err = db.c.QueryRow(queryChatmate, c.ConversationID, id).Scan(&u.UserID, &u.Username, &u.Image); err != nil {
			return nil, err
		}

		pc.User = &u
		c.Private = &pc
		res = append(res, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *appdbimpl) AddPrivateChat(senderID uuid.UUID, mpb MessagePrivateBody) (Conversation, error) {

	var res Conversation

	const queryMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, message, image, hasImage)
		VALUES (?,?,?,?,?, ?);
	`

	const queryAddConversation = `
		INSERT INTO Conversation(conversationID, conversationType) VALUES
		(?,'private');
	`

	const queryAddMembers = `
		INSERT INTO Members(userID, conversationID) VALUES
		(?,?);

	`

	const queryLatestMessage = `
	SELECT m.timestamp, m.message, m.image
	FROM Message m 
	WHERE m.conversationID = ?
	ORDER BY m.timestamp DESC
	LIMIT 1;
	`

	const queryReceiver = `
		SELECT u.userID, u.username, u.image
		FROM User u
		WHERE u.userID = ?;
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

	convID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	// Add Conversation

	if _, err = tx.Exec(queryAddConversation, convID); err != nil {
		return res, err
	}

	// Add Members
	mems := [2]uuid.UUID{senderID, mpb.ReceiverID}
	for _, id := range mems {

		if _, err = tx.Exec(queryAddMembers, id, convID); err != nil {
			return res, err
		}

	}

	messageID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	// Add Message
	var hasImage bool
	if mpb.Image != nil {
		hasImage = true
	}
	if _, err = tx.Exec(queryMessage, messageID, senderID, convID, mpb.Message, mpb.Image, hasImage); err != nil {
		return res, err
	}

	// Query Receiver
	var u User
	if err = tx.QueryRow(queryReceiver, mpb.ReceiverID).Scan(&u.UserID, &u.Username, &u.Image); err != nil {
		return res, err
	}

	// Query Latest Message
	var lm LatestMessage
	if err = tx.QueryRow(queryLatestMessage, convID).Scan(&lm.Timestamp, &lm.Message, &lm.Image); err != nil {
		return res, err
	}

	if err = tx.Commit(); err != nil {
		return res, err
	}

	pc := &PrivateConversation{
		User: &u,
	}

	res.ConversationID = convID
	res.Type = "private"
	res.Private = pc
	res.LatestMessage = &lm

	return res, nil

}
