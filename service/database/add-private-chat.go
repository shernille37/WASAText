package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type MessagePrivateBody struct {
	ReceiverID  uuid.UUID
	Message     string
	MessageType string
}

func (db *appdbimpl) AddPrivateChat(senderID uuid.UUID, mpb MessagePrivateBody) (PrivateConversation, error) {

	var res PrivateConversation

	const queryMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, messageType, messageStatus, timeDelivered, message, image)
		VALUES (?,?,?,?,'sent', "", ?, ?);
	`

	const queryAddConversation = `
		INSERT INTO Conversation(conversationID, conversationType, groupName, groupImage) VALUES
		(?,?,"","");
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

	const queryCheckConversation = `
		SELECT EXISTS( SELECT c.conversationID 
		FROM Conversation c, Members m, Message mess
		WHERE mess.senderID = ? AND mess.conversationID = c.conversationID 
		AND m.conversationID = c.conversationID AND m.userID = ?);
	`

	const queryLatestMessage = `
	SELECT m.hasImage, m.timestamp, COALESCE(m.message, '') 
	FROM Message m 
	WHERE m.conversationID = ?
	ORDER BY m.timestamp DESC
	LIMIT 1;
	`

	if senderID == mpb.ReceiverID {
		return res, fmt.Errorf("cannot have conversation to yourself")
	}

	// No user found -- Register the User to the system
	var u User
	if err := db.c.QueryRow(queryCheckUser, mpb.ReceiverID).Scan(&u.UserID, &u.Name, &u.Image); errors.Is(err, sql.ErrNoRows) {
		return res, fmt.Errorf("user doesn't exists")

	}

	var existingConvID bool
	if err := db.c.QueryRow(queryCheckConversation, senderID, mpb.ReceiverID).Scan(&existingConvID); err != nil {
		return res, err
	}

	if existingConvID {
		return res, fmt.Errorf("conversation already exists")
	}

	convID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	messageID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	// Add Conversation
	var mess string
	var image string

	if mpb.MessageType == "image" {
		image = mpb.Message
	} else {

		mess = mpb.Message
	}

	_, err = db.c.Exec(queryAddConversation, convID, "personal")

	if err != nil {
		return res, err
	}

	// Add Members
	mems := [2]uuid.UUID{senderID, mpb.ReceiverID}
	for _, id := range mems {

		_, err = db.c.Exec(queryAddMembers, id, convID)

		if err != nil {
			return res, err
		}
	}

	// Add Message
	_, err = db.c.Exec(queryMessage, messageID, senderID, convID, mpb.MessageType, mess, image)

	if err != nil {
		return res, err
	}

	var lm LatestMessage
	if err = db.c.QueryRow(queryLatestMessage, convID).Scan(&lm.HasImage, &lm.Timestamp, &lm.Message); err != nil {
		return res, err
	}

	res.ConversationID = convID
	res.User = &u
	res.LatestMessage = &lm

	return res, nil

}
