package database

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type MessageGroupBody struct {
	GroupName  string
	GroupImage *string
	Members    []uuid.UUID
	Message    string
	Image      *string
}

type GroupMemberBody struct {
	Members []uuid.UUID
}

func (db *appdbimpl) ListGroupConversation(id uuid.UUID) ([]Conversation, error) {

	var res []Conversation

	const queryGroupConversation = `SELECT c.conversationID, c.conversationType, c.groupName, c.groupImage 
	FROM Conversation c, Members m, User u
	WHERE c.conversationType = 'group' AND u.userID = ? AND c.conversationID = m.conversationID AND m.userID = u.userID;
	`

	rows, err := db.c.Query(queryGroupConversation, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c Conversation
		var gc GroupConversation

		if err = rows.Scan(&c.ConversationID, &c.Type, &gc.GroupName, &gc.GroupImage); err != nil {
			return nil, err
		}
		c.Group = &gc
		res = append(res, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil

}

func (db *appdbimpl) AddGroupConversation(senderID uuid.UUID, mgb MessageGroupBody) (Conversation, error) {

	var res Conversation

	const queryAddMessage = `
		INSERT INTO Message(messageID, senderID, conversationID, hasImage, message, image)
		VALUES (?,?,?,?,?,?);
	`

	const queryAddConversation = `
		INSERT INTO Conversation(conversationID, conversationType, groupName, groupImage) VALUES
		(?,'group',?,?);
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

	tx, err := db.c.Begin()
	if err != nil {
		return res, err
	}

	// In case of error --> Rollback
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed %v\n", rbErr)
			}
		}
	}()

	convID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	messageID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	// Add conversation
	if _, err = tx.Exec(queryAddConversation, convID, mgb.GroupName, mgb.GroupImage); err != nil {
		return res, err
	}

	// Add members
	// Add the creator of the group
	mgb.Members = append(mgb.Members, senderID)
	for _, userID := range mgb.Members {
		if _, err = tx.Exec(queryAddMembers, userID, convID); err != nil {
			return res, err
		}
	}

	// Add Message
	var hasImage bool
	if mgb.Image != nil {
		hasImage = true
	}

	if _, err = tx.Exec(queryAddMessage, messageID, senderID, convID, hasImage, mgb.Message, mgb.Image); err != nil {
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

	gc := &GroupConversation{
		GroupName:  mgb.GroupName,
		GroupImage: mgb.GroupImage,
	}

	res.ConversationID = convID
	res.Type = "group"
	res.Group = gc
	res.LatestMessage = &lm

	return res, nil
}

func (db *appdbimpl) UpdateGroupName(conversationID uuid.UUID, newGroupName string) error {

	const queryUpdate = `
		UPDATE Conversation
		SET groupName = ?
		WHERE conversationID = ? AND conversationType = 'group';
	`

	res, err := db.c.Exec(queryUpdate, newGroupName, conversationID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Conversation does not exist")
	}

	return nil
}

func (db *appdbimpl) UpdateGroupImage(conversationID uuid.UUID, newGroupPhoto string) error {

	const queryUpdate = `
		UPDATE Conversation
		SET groupImage = ?
		WHERE conversationID = ? AND conversationType = 'group';
	`

	res, err := db.c.Exec(queryUpdate, newGroupPhoto, conversationID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Conversation does not exist")
	}

	return nil
}

func (db *appdbimpl) ListGroupMembers(conversationID uuid.UUID) ([]User, error) {
	var res []User

	const queryMembers = `
		SELECT u.userID, u.username, u.image FROM 
		Members m, User u
		WHERE m.conversationID = ? AND m.userID = u.userID;
	`

	rows, err := db.c.Query(queryMembers, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		if err = rows.Scan(&u.UserID, &u.Username, &u.Image); err != nil {
			return nil, err
		}

		res = append(res, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *appdbimpl) AddGroupMembers(conversationID uuid.UUID, gmb GroupMemberBody) ([]User, error) {

	var res []User

	const queryAddMembers = `
		INSERT INTO Members(userID, conversationID)
		VALUES (?, ?);
	`

	const queryResponse = `
	SELECT u.userID, u.username, u.image 
	FROM User u
	WHERE u.userID = ?;
	`

	tx, err := db.c.Begin()

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Rollback failed %v\n", rbErr)
			}
		}
	}()

	for _, member := range gmb.Members {

		_, err := tx.Exec(queryAddMembers, member, conversationID)

		if err != nil {
			return nil, err
		}
	}

	for _, member := range gmb.Members {
		var u User
		if err := tx.QueryRow(queryResponse, member).Scan(&u.UserID, &u.Username, &u.Image); err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	if err = tx.Commit(); err != nil {
		return res, err
	}

	return res, nil

}

func (db *appdbimpl) LeaveGroupConversation(userID uuid.UUID, conversationID uuid.UUID) error {

	const queryDelete = `
		DELETE FROM Members
		WHERE conversationID = ? AND userID = ?;
	`

	res, err := db.c.Exec(queryDelete, conversationID, userID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Conversation does not exist")
	}

	return nil
}
