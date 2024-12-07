package database

import "github.com/gofrs/uuid"

type User struct {
	UserID uuid.UUID
	Name   string
	Image  string
}

type PrivateConversation struct {
	ConversationID uuid.UUID
	User           *User
	LatestMessage  *LatestMessage
}

type GroupConversation struct {
	ConversationID uuid.UUID
	GroupName      string
	GroupImage     string
	LatestMessage  *LatestMessage
}

type LatestMessage struct {
	MessageType string
	Message     string
	Timestamp   string
}

type Conversation struct {
	Type    string
	Private *PrivateConversation
	Group   *GroupConversation
	Members []User
}

type Message struct {
	MessageID      uuid.UUID
	SenderID       uuid.UUID
	ConversationID uuid.UUID
	ReplyMessageID *uuid.UUID
	ReplyMessage   *string
	Timestamp      string
	HasImage       bool
	MessageType    string
	MessageStatus  string
	TimeDelivered  string
	Message        string
	Image          *string
}
