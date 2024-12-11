package database

import "github.com/gofrs/uuid"

type ForwardMessageBody struct {
	Source      uuid.UUID `json:"source"`
	Destination uuid.UUID `json:"destination"`
}

func (db *appdbimpl) ForwardMessage(senderID uuid.UUID, messageID uuid.UUID, fmb ForwardMessageBody) error {

	const queryMessage = `
		SELECT m.message, m.hasImage, COALESCE(m.image, '')
		FROM Message m
		WHERE m.messageID = ? AND m.conversationID = ?; 
	`

	const queryForward = `
		INSERT INTO Message(messageID, senderID, conversationID, messageType, hasImage, message, image)
		VALUES (?,?,?,?,?,?,?);
	`

	var sourceMessage Message

	if err := db.c.QueryRow(queryMessage, messageID, fmb.Source).Scan(&sourceMessage.Message, &sourceMessage.HasImage, &sourceMessage.Image); err != nil {
		return err
	}

	newMessageID, err := uuid.NewV4()

	if err != nil {
		return err
	}

	sourceMessage.MessageType = "forward"

	if !sourceMessage.HasImage {
		sourceMessage.Image = nil
	}

	_, err = db.c.Exec(queryForward, newMessageID, senderID, fmb.Destination, sourceMessage.MessageType, sourceMessage.HasImage, sourceMessage.Message, sourceMessage.Image)

	if err != nil {
		return err
	}

	return nil

}
