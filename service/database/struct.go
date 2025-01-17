package database

import "github.com/gofrs/uuid"

type User struct {
	UserID   uuid.UUID
	Username string
	Image    *string
}

type PrivateConversation struct {
	User *User
}

type GroupConversation struct {
	GroupName  string
	GroupImage *string
}

type LatestMessage struct {
	Message   string
	Image     *string
	Timestamp string
}

type Conversation struct {
	ConversationID uuid.UUID
	Type           string
	Private        *PrivateConversation
	Group          *GroupConversation
	LatestMessage  *LatestMessage
	Members        []User
}

type Message struct {
	MessageID          uuid.UUID
	SenderID           uuid.UUID
	SenderName         string
	ConversationID     uuid.UUID
	ReplyMessageID     *uuid.UUID
	ReplyMessage       *string
	ReplyRecipientName *string
	Timestamp          string
	HasImage           bool
	MessageType        string
	MessageStatus      string
	TimeDelivered      *string
	TimeRead           *string
	Message            string
	Image              *string
	Reactions          []Reaction
}

type Reader struct {
	User      *User
	Timestamp string
}

type Reaction struct {
	ReactionID *uuid.UUID
	Unicode    string
	Reactor    *User
	Count      *int
}
