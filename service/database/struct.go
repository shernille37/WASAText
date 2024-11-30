package database

import "github.com/gofrs/uuid"

type PrivateConversation struct {
	ConversationID uuid.UUID
	UserID uuid.UUID
	Name string
	Image string
	LatestMessage *LatestMessage

}

type GroupConversation struct {
	ConversationID uuid.UUID
	GroupName string
	GroupImage string
	LatestMessage *LatestMessage
}


type LatestMessage struct {
	MessageType string
	Message string
	Timestamp string
}

type Conversation struct {
	Type string
	Private *PrivateConversation
	Group *GroupConversation
}
