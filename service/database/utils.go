package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/shernille37/WASAText/service/api/constants"
)

func (db *appdbimpl) GetUserByID(id string) (User, error) {

	var res User
	const queryUser = `
		SELECT u.userID, u.username, u.image
		FROM User u
		WHERE u.userID = ?;
	`

	// No User found
	if err := db.c.QueryRow(queryUser, id).Scan(&res.UserID, &res.Username, &res.Image); errors.Is(err, sql.ErrNoRows) {
		return res, err
	}

	return res, nil
}

func (db *appdbimpl) CheckConversationMembership(id uuid.UUID, members []uuid.UUID) error {

	const queryMembership = `
		SELECT EXISTS(SELECT 1 FROM Members mem
		WHERE mem.conversationID = ? AND mem.userID = ?);
	`

	var res bool
	for _, member := range members {
		if _ = db.c.QueryRow(queryMembership, id, member).Scan(&res); !res {
			return errors.New(constants.ERROR_USER_MEMBERSHIP)
		}
	}

	return nil
}

func (db *appdbimpl) CheckMessageMembership(conversationID uuid.UUID, messageID uuid.UUID) error {

	const queryMembership = `
		SELECT EXISTS(SELECT 1 FROM Message mess
		WHERE mess.conversationID = ? AND mess.messageID = ?);
	`

	var res bool

	if _ = db.c.QueryRow(queryMembership, conversationID, messageID).Scan(&res); !res {
		return errors.New(constants.ERROR_MESSAGE_MEMBERSHIP)
	}

	return nil
}

func (db *appdbimpl) CheckExistingConversation(senderID uuid.UUID, receiverID uuid.UUID) error {
	const queryCheckConversation = `
	SELECT EXISTS (SELECT 1 
	FROM Conversation c, Conversation c1, Members m1, Members m2
	WHERE m1.userID = ? AND m2.userID = ? AND 
	m1.conversationID = c.conversationID AND 
	m2.conversationID = c1.conversationID AND
	c.conversationID = c1.conversationID AND c.conversationType = 'personal' AND c1.conversationType = 'personal');
`

	var existingConvID bool
	if err := db.c.QueryRow(queryCheckConversation, senderID, receiverID).Scan(&existingConvID); err != nil {
		return err
	}

	if existingConvID {
		return errors.New(constants.CONVERSATION_ALREADY_EXISTS)
	}

	return nil
}

func (db *appdbimpl) CheckMessageOwnership(userID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) error {

	const queryMembership = `
		SELECT EXISTS(SELECT 1 FROM Message mess
		WHERE mess.conversationID = ? AND mess.messageID = ? AND mess.senderID = ?);
	`

	var res bool

	if _ = db.c.QueryRow(queryMembership, conversationID, messageID, userID).Scan(&res); !res {
		return errors.New(constants.UNAUTHORIZED)
	}

	return nil
}

func (db *appdbimpl) CheckReactionOwnership(userID uuid.UUID, reactionID uuid.UUID) error {

	const queryMembership = `
		SELECT EXISTS(SELECT 1 FROM Reaction r
		WHERE r.reactionID = ? AND r.userID = ?);
	`

	var res bool

	if _ = db.c.QueryRow(queryMembership, reactionID, userID).Scan(&res); !res {
		return errors.New(constants.UNAUTHORIZED)
	}

	return nil
}
