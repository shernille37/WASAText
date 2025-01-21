package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type ReactionBody struct {
	Unicode string `json:"unicode"`
}

func (db *appdbimpl) ListEmojis() ([]string, error) {

	var res []string

	const queryEmojis = `
		SELECT unicode 
		FROM Emoji;
	`

	rows, err := db.c.Query(queryEmojis)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var emoji string

		if err = rows.Scan(&emoji); err != nil {
			return nil, err
		}

		res = append(res, emoji)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *appdbimpl) ListReactions(conversationID uuid.UUID, messageID uuid.UUID) ([]Reaction, error) {

	var res []Reaction

	const queryReactions = `
		SELECT r.reactionID, r.emoji, u.userID, u.username, u.image 
		FROM Reaction r, User u, Message m
		WHERE r.messageID = ? AND r.messageID = m.messageID AND 
		r.userID = u.userID AND m.conversationID = ?;
	`

	rows, err := db.c.Query(queryReactions, messageID, conversationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var r Reaction
		var u User

		if err = rows.Scan(&r.ReactionID, &r.Unicode, &u.UserID, &u.Username, &u.Image); err != nil {
			return nil, err
		}

		r.Reactor = &u

		res = append(res, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil

}

func (db *appdbimpl) AddReaction(userID uuid.UUID, messageID uuid.UUID, rb ReactionBody) (Reaction, error) {

	var res Reaction

	const queryAddReaction = `
		INSERT INTO Reaction(reactionID, userID, messageID, emoji) VALUES 
		(?, ?, ?, ?);
	`

	const queryResponse = `
		SELECT r.reactionID, r.emoji, u.userID, u.username, u.image 
		FROM Reaction r, User u
		WHERE r.messageID = ? AND r.reactionID = ? AND 
		r.userID = u.userID;
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

	// Add reaction
	reactionID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	if _, err := tx.Exec(queryAddReaction, reactionID, userID, messageID, rb.Unicode); err != nil {
		return res, err
	}

	var u User
	if err = tx.QueryRow(queryResponse, messageID, reactionID).Scan(&res.ReactionID, &res.Unicode, &u.UserID, &u.Username, &u.Image); err != nil {
		return res, err
	}

	if err = tx.Commit(); err != nil {
		return res, err
	}

	res.Reactor = &u

	return res, nil

}

func (db *appdbimpl) DeleteReaction(reactionID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) error {

	var res sql.Result

	const queryDeleteReaction = `
		DELETE FROM Reaction WHERE
		reactionID = ? AND messageID = ? AND 
		messageID IN (
			SELECT m.messageID
			FROM Message m
			WHERE m.conversationID = ? 
		);
	`

	res, err := db.c.Exec(queryDeleteReaction, reactionID, messageID, conversationID)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Reaction does not exist")
	}

	return nil

}
