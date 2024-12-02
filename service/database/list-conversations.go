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