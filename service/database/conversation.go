package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListConversation(id uuid.UUID) ([]Conversation, error) {

	var res []Conversation

	personalConversations, err1 := db.ListPrivateConversation(id)
	groupConversations, err2 := db.ListGroupConversation(id)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	for _, pc := range personalConversations {
		var c Conversation
		tmp := pc
		c.Type = "personal"
		c.Private = &tmp
		res = append(res, c)
	}

	for _, gc := range groupConversations {
		var c Conversation
		tmp := gc
		c.Type = "group"
		c.Group = &tmp
		res = append(res, c)
	}

	return res, nil

}

func (db *appdbimpl) GetConversation(id uuid.UUID, conversationID uuid.UUID) (Conversation, error) {
	var res Conversation
	var members []User

	const queryConversation = `
		SELECT c.conversationID, c.conversationType, c.groupName, c.groupImage
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
		SELECT m.timestamp, m.message, m.image
		FROM Message m 
		WHERE m.conversationID = ?
		ORDER BY m.timestamp DESC
		LIMIT 1;
	`

	const queryMembers = `
		SELECT u.userID, u.username, u.image FROM 
		Members m, User u
		WHERE m.conversationID = ? AND m.userID = u.userID;
	`

	var convID uuid.UUID
	var conversationType string
	var groupName *string
	var groupImage *string

	if err := db.c.QueryRow(queryConversation, conversationID).Scan(&convID, &conversationType, &groupName, &groupImage); err != nil {
		return res, err
	}

	res.Type = conversationType

	if conversationType == "personal" {
		var pc PrivateConversation
		var u User
		if err := db.c.QueryRow(queryPrivateConversation, conversationID, id).Scan(&pc.ConversationID, &u.UserID, &u.Name, &u.Image); err != nil {
			return res, err
		}

		pc.User = &u
		res.Private = &pc
	} else {
		var gc GroupConversation
		gc.ConversationID = convID
		gc.GroupName = *groupName
		gc.GroupImage = groupImage

		res.Group = &gc
	}

	var lm LatestMessage
	if err := db.c.QueryRow(queryLatestMessage, conversationID).Scan(&lm.Timestamp, &lm.Message, &lm.Image); err != nil {
		return res, err
	}

	if res.Private != nil {
		res.Private.LatestMessage = &lm
	} else {
		res.Group.LatestMessage = &lm
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

	if err = rows.Err(); err != nil {
		return res, err
	}

	res.Members = members

	return res, nil
}