package database

func (db *appdbimpl) ListPrivateConversation(id string) ([]PrivateConversation, error) {

	var res []PrivateConversation

	const queryPrivateConversation = `SELECT c.id, u.id, u.username, u.image
	FROM Conversation c, Members m, User u
	WHERE u.id = ? AND c.id = m.conversationID AND m.userID = ?
	`
	const queryLatestMessage = `
		SELECT	m.messageType, m.timestamp, m.message FROM Message m 
		WHERE m.conversationID = ?
		ORDER BY m.timestamp
		LIMIT 1
	`

	rows, err := db.c.Query(queryPrivateConversation, id, id)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var pc PrivateConversation
		var lm LatestMessage
		

		err = rows.Scan(&pc.ConversationID, &pc.UserID, &pc.Name, &pc.Image)
		if err != nil {
			return nil, err
		}

		rowLatestMessage, err := db.c.Query(queryLatestMessage, pc.ConversationID)
		if err != nil {
			return nil, err
		}
		defer func() { _ = rowLatestMessage.Close() }()

		for rowLatestMessage.Next() {
			err = rowLatestMessage.Scan(&lm.MessageType, &lm.Timestamp, &lm.Message)
			if err != nil {
				return nil, err
			}
		}

		pc.LatestMessage = lm

		res = append(res, pc)
	}

	return res, nil

}