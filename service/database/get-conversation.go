package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) GetConversation(id uuid.UUID, conversationID uuid.UUID) (Conversation, error) {
	var res Conversation
	var members []User

	
	const queryConversation = `
		SELECT c.conversationID, c.conversationType, COALESCE(c.groupName, ''), COALESCE(c.groupImage, '')
		FROM Conversation c
		WHERE c.conversationID = ?;
	`
	
	const queryPrivateConversation = `
	SELECT c.conversationID, u.userID, u.username, u.image
	FROM Conversation c, Members m, User u
	WHERE c.conversationType = 'personal' AND c.conversationID = ? AND c.conversationID = m.conversationID AND m.userID <> ?
	AND m.userID = u.userID;

	`

	const queryLatestMessage = `
		SELECT m.messageType, m.timestamp, COALESCE(m.message, '') 
		FROM Message m 
		WHERE m.conversationID = ?
		ORDER BY m.timestamp DESC
		LIMIT 1;
	`

	const queryMembers = `
		SELECT u.userID, u.username, u.image FROM 
		Members m, User u
		WHERE m.conversationID = ? AND m.userID = u.userID AND m.userID <> ?;
	`

	var convID uuid.UUID
	var convType string
	var groupName string
	var groupImage string

	if err := db.c.QueryRow(queryConversation, conversationID).Scan(&convID, &convType, &groupName, &groupImage); err != nil {
		return res, err
	}


	res.Type = convType

	if convType == "personal" {
		var pc PrivateConversation
		var u User
		if err := db.c.QueryRow(queryPrivateConversation, conversationID, id).Scan(&pc.ConversationID, &u.UserID, &u.Name, &u.Image); err != nil {
			return res, err
		}


		var lm LatestMessage
		if err := db.c.QueryRow(queryLatestMessage, conversationID).Scan(&lm.MessageType, &lm.Timestamp, &lm.Message); err != nil {
			return res, err
		}

		pc.User = &u
		pc.LatestMessage = &lm
		res.Private = &pc
	} else {
		var gc GroupConversation
		gc.ConversationID = convID
		gc.GroupName = groupName
		gc.GroupImage = groupImage


		var lm LatestMessage
		if err := db.c.QueryRow(queryLatestMessage, conversationID).Scan(&lm.MessageType, &lm.Timestamp, &lm.Message); err != nil {
			return res, err
		}

		gc.LatestMessage = &lm
		res.Group = &gc
	}

	rows, err := db.c.Query(queryMembers, conversationID, id)
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User

		if err = rows.Scan(&u.UserID, &u.Name, &u.Image); err != nil {
			return res, err
		}

		members = append(members, u)
	}

	res.Members = members
	return res, nil
}
