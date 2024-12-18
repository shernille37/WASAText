package database

import "github.com/gofrs/uuid"

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

		if err = rows.Scan(&r.ReactionID, &r.Unicode, &u.UserID, &u.Name, &u.Image); err != nil {
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
