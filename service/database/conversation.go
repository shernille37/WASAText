package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ListConversation(id uuid.UUID) ([]Conversation, error) {

	var res []Conversation

	const queryLatestMessage = `
	SELECT m.timestamp, m.message, m.image
	FROM Message m 
	WHERE m.conversationID = ?
	ORDER BY m.timestamp DESC
	LIMIT 1;
	`

	personalConversations, err1 := db.ListPrivateConversation(id)
	groupConversations, err2 := db.ListGroupConversation(id)

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	res = append(res, personalConversations...)
	res = append(res, groupConversations...)

	for idx := range res {
		lm := &LatestMessage{}
		err := db.c.QueryRow(queryLatestMessage, res[idx].ConversationID).Scan(&lm.Timestamp, &lm.Message, &lm.Image)
		if errors.Is(err, sql.ErrNoRows) {
			res[idx].LatestMessage = nil
		} else {
			res[idx].LatestMessage = lm
		}
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
	SELECT u.userID, u.username, u.image
	FROM Conversation c, Members m, User u
	WHERE c.conversationType = 'private' AND c.conversationID = ? AND c.conversationID = m.conversationID AND m.userID <> ?
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

	var groupName *string
	var groupImage *string

	if err := db.c.QueryRow(queryConversation, conversationID).Scan(&res.ConversationID, &res.Type, &groupName, &groupImage); err != nil {
		return res, err
	}

	if res.Type == "private" {
		var pc PrivateConversation
		var u User
		if err := db.c.QueryRow(queryPrivateConversation, conversationID, id).Scan(&u.UserID, &u.Username, &u.Image); err != nil {
			return res, err
		}

		pc.User = &u
		res.Private = &pc
	} else {
		var gc GroupConversation

		gc.GroupName = *groupName
		gc.GroupImage = groupImage

		res.Group = &gc
	}

	// Query Latest Message
	var lm LatestMessage
	if err := db.c.QueryRow(queryLatestMessage, conversationID).Scan(&lm.Timestamp, &lm.Message, &lm.Image); errors.Is(err, sql.ErrNoRows) {
		res.LatestMessage = nil
	} else {
		res.LatestMessage = &lm
	}

	// Query Members
	rows, err := db.c.Query(queryMembers, conversationID, id)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User

		if err = rows.Scan(&u.UserID, &u.Username, &u.Image); err != nil {
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
