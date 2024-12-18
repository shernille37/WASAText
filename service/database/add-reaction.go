package database

import "github.com/gofrs/uuid"

type ReactionBody struct {
	Unicode string `json:"unicode"`
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

	// Add Reaction

	reactionID, err := uuid.NewV4()
	if err != nil {
		return res, err
	}

	if _, err := db.c.Exec(queryAddReaction, reactionID, userID, messageID, rb.Unicode); err != nil {
		return res, err
	}

	var u User
	if err = db.c.QueryRow(queryResponse, messageID, reactionID).Scan(&res.ReactionID, &res.Unicode, &u.UserID, &u.Name, &u.Image); err != nil {
		return res, err
	}

	res.Reactor = &u

	return res, nil

}
