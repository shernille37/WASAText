package database

import (
	"github.com/gofrs/uuid"
)


type MessagePrivateBody struct {
	ReceiverID uuid.UUID
	Message    string
	Image      *string
}

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
		SELECT m.timestamp, m.message, m.image
		FROM Message m 
		WHERE m.conversationID = ?
		ORDER BY m.timestamp DESC
		LIMIT 1;
	`

	rows, err := db.c.Query(queryPrivateConversation, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		var pc PrivateConversation
		var lm LatestMessage

		if err = rows.Scan(&pc.ConversationID); err != nil {
			return nil, err
		}

		// Fetch chatmate
		if err = db.c.QueryRow(queryChatmate, pc.ConversationID, id).Scan(&u.UserID, &u.Name, &u.Image); err != nil {
			return nil, err
		}

		// Fetch the latest message
		if err = db.c.QueryRow(queryLatestMessage, pc.ConversationID).Scan(&lm.Timestamp, &lm.Message, &lm.Image); err != nil {
			return nil, err
		}

		pc.User = &u
		pc.LatestMessage = &lm
		res = append(res, pc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *appdbimpl) AddPrivateChat(senderID uuid.UUID, mpb MessagePrivateBody) (PrivateConversation, error) {

	var res PrivateConversation

	const queryMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, message, image)
		VALUES (?,?,?,?,?);
	`

	const queryAddConversation = `
		INSERT INTO Conversation(conversationID, conversationType) VALUES
		(?,'personal');
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
			tx.Rollback()
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
	if _, err = tx.Exec(queryMessage, messageID, senderID, convID, mpb.Message, mpb.Image); err != nil {
		return res, err
	}

	var u User
	if err = tx.QueryRow(queryReceiver, mpb.ReceiverID).Scan(&u.UserID, &u.Name, &u.Image); err != nil {
		return res, err
	}

	var lm LatestMessage
	if err = tx.QueryRow(queryLatestMessage, convID).Scan(&lm.Timestamp, &lm.Message, &lm.Image); err != nil {
		return res, err
	}

	if err = tx.Commit(); err != nil {
		return res, err
	}

	res.ConversationID = convID
	res.User = &u
	res.LatestMessage = &lm

	return res, nil

}