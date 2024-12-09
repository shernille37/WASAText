package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type MessageGroupBody struct {
	GroupName   string
	GroupImage  string
	MessageType string
	Message     string
	Members     []uuid.UUID
}

func (db *appdbimpl) AddGroupChat(senderID uuid.UUID, mgb MessageGroupBody) (GroupConversation, error) {

	var res GroupConversation

	const queryMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, messageType, message, image)
		VALUES (?,?,?,?, ?, ?);
	`

	const queryAddConversation = `
		INSERT INTO Conversation(conversationID, conversationType, groupName, groupImage) VALUES
		(?,?,?,?);
	`

	const queryAddMembers = `
		INSERT INTO Members(userID, conversationID) VALUES
		(?,?);
	`

	const queryCheckUser = `
		SELECT u.userID, u.username, u.image 
		FROM User u
		WHERE u.userID = ?;
	`

	const queryLatestMessage = `
	SELECT m.hasImage, m.timestamp, COALESCE(m.message, '') 
	FROM Message m 
	WHERE m.conversationID = ?
	ORDER BY m.timestamp DESC
	LIMIT 1;
	`

	var u_id uuid.UUID
	mgb.Members = append(mgb.Members, senderID)
	// Check for user existence in the system
	for _, userID := range mgb.Members {

		if err := db.c.QueryRow(queryCheckUser, userID).Scan(&u_id); errors.Is(err, sql.ErrNoRows) {
			return res, fmt.Errorf("user doesn't exists")
		}

	}

	convID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	messageID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	// Add conversation
	if _, err = db.c.Exec(queryAddConversation, convID, "group", mgb.GroupName, mgb.GroupImage); err != nil {
		return res, err
	}

	// Add members
	for _, userID := range mgb.Members {
		if _, err = db.c.Exec(queryAddMembers, userID, convID); err != nil {
			return res, err
		}
	}

	// Add Message
	var mess *string
	var image *string

	if mgb.MessageType == "image" {
		image = &mgb.Message
	} else {
		mess = &mgb.Message
	}

	_, err = db.c.Exec(queryMessage, messageID, senderID, convID, mgb.MessageType, mess, image)

	if err != nil {
		return res, err
	}

	var lm LatestMessage
	err = db.c.QueryRow(queryLatestMessage, convID).Scan(&lm.HasImage, &lm.Timestamp, &lm.Message)

	if err != nil {
		return res, err
	}

	res.ConversationID = convID
	res.GroupName = mgb.GroupName
	res.GroupImage = mgb.GroupImage
	res.LatestMessage = &lm

	return res, nil

}
