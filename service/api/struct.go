package api

import (
	"github.com/gofrs/uuid"
	"github.com/shernille37/WASAText/service/database"
)

type User struct {
	UserID uuid.UUID `json:"userID"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
}

type PrivateConversation struct {
	ConversationID uuid.UUID      `json:"conversationID"`
	User           *User          `json:"user"`
	LatestMessage  *LatestMessage `json:"latestMessage"`
}

type GroupConversation struct {
	ConversationID uuid.UUID      `json:"conversationID"`
	GroupName      string         `json:"groupName"`
	GroupImage     string         `json:"groupImage"`
	LatestMessage  *LatestMessage `json:"latestMessage"`
}

type LatestMessage struct {
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
}

type Conversation struct {
	Type    string               `json:"type"`
	Private *PrivateConversation `json:"private"`
	Group   *GroupConversation   `json:"group"`
	Members []User               `json:"members"`
}

type Message struct {
	MessageID      uuid.UUID  `json:"messageID"`
	SenderID       uuid.UUID  `json:"senderID"`
	ConversationID uuid.UUID  `json:"conversationID"`
	ReplyMessageID *uuid.UUID `json:"replyMessageID"`
	ReplyMessage   *string    `json:"replyMessage"`
	Timestamp      string     `json:"timestamp"`
	HasImage       bool       `json:"hasImage"`
	MessageType    string     `json:"messageType"`
	MessageStatus  string     `json:"messageStatus"`
	TimeDelivered  string     `json:"timeDelivered"`
	Message        string     `json:"message"`
	Image          *string    `json:"image"`
}

func (m *Message) ToDatabase() database.Message {
	return database.Message{
		MessageID:      m.MessageID,
		SenderID:       m.SenderID,
		ConversationID: m.ConversationID,
		ReplyMessageID: m.ReplyMessageID,
		ReplyMessage:   m.ReplyMessage,
		Timestamp:      m.Timestamp,
		HasImage:       m.HasImage,
		MessageType:    m.MessageType,
		MessageStatus:  m.MessageStatus,
		TimeDelivered:  m.TimeDelivered,
		Message:        m.Message,
		Image:          m.Image,
	}
}

func (m *Message) FromDatabase(mess database.Message) {
	m.MessageID = mess.MessageID
	m.SenderID = mess.SenderID
	m.ConversationID = mess.ConversationID
	m.ReplyMessageID = mess.ReplyMessageID
	m.ReplyMessage = mess.ReplyMessage
	m.Timestamp = mess.Timestamp
	m.HasImage = mess.HasImage
	m.MessageType = mess.MessageType
	m.MessageStatus = mess.MessageStatus
	m.TimeDelivered = mess.TimeDelivered
	m.Message = mess.Message
	m.Image = mess.Image
}

func (u *User) FromDatabase(user database.User) {
	u.UserID = user.UserID
	u.Name = user.Name
	u.Image = user.Image
}

func (u *User) ToDatabase() database.User {
	return database.User{
		UserID: u.UserID,
		Name:   u.Name,
		Image:  u.Image,
	}
}

func (c *Conversation) ToDatabase() database.Conversation {
	var private *database.PrivateConversation
	if c.Private != nil {
		private = &database.PrivateConversation{
			ConversationID: c.Private.ConversationID,
			User:           (*database.User)(c.Private.User),
			LatestMessage:  (*database.LatestMessage)(c.Private.LatestMessage),
		}
	}

	var group *database.GroupConversation
	if c.Group != nil {
		group = &database.GroupConversation{
			ConversationID: c.Group.ConversationID,
			GroupName:      c.Group.GroupName,
			GroupImage:     c.Group.GroupImage,
			LatestMessage:  (*database.LatestMessage)(c.Group.LatestMessage),
		}
	}

	var members []database.User

	for _, member := range c.Members {
		members = append(members, database.User{
			UserID: member.UserID,
			Name:   member.Name,
			Image:  member.Image,
		})
	}

	return database.Conversation{
		Type:    c.Type,
		Private: private,
		Group:   group,
		Members: members,
	}
}

func (c *Conversation) FromDatabase(conv database.Conversation) {
	c.Type = conv.Type

	if conv.Private != nil {
		c.Private = &PrivateConversation{
			ConversationID: conv.Private.ConversationID,
			User:           (*User)(conv.Private.User),
			LatestMessage:  (*LatestMessage)(conv.Private.LatestMessage),
		}
	}

	if conv.Group != nil {
		c.Group = &GroupConversation{
			ConversationID: conv.Group.ConversationID,
			GroupName:      conv.Group.GroupName,
			GroupImage:     conv.Group.GroupImage,
			LatestMessage:  (*LatestMessage)(conv.Group.LatestMessage),
		}

	}

	var members = make([]User, len(conv.Members))
	for idx := range members {
		members[idx] = User{
			UserID: conv.Members[idx].UserID,
			Name:   conv.Members[idx].Name,
			Image:  conv.Members[idx].Image,
		}
	}

	c.Members = members

}

func (pc *PrivateConversation) FromDatabase(priv database.PrivateConversation) {
	pc.ConversationID = priv.ConversationID
	pc.User = (*User)(priv.User)
	pc.LatestMessage = (*LatestMessage)(priv.LatestMessage)
}

func (pc *PrivateConversation) ToDatabase() database.PrivateConversation {
	return database.PrivateConversation{
		ConversationID: pc.ConversationID,
		User:           (*database.User)(pc.User),
		LatestMessage:  (*database.LatestMessage)(pc.LatestMessage),
	}
}

func (gc *GroupConversation) FromDatabase(group database.GroupConversation) {
	gc.ConversationID = group.ConversationID
	gc.GroupName = group.GroupName
	gc.GroupImage = group.GroupImage
	gc.LatestMessage = (*LatestMessage)(group.LatestMessage)
}

func (gc *GroupConversation) ToDatabase() database.GroupConversation {
	return database.GroupConversation{
		ConversationID: gc.ConversationID,
		GroupName:      gc.GroupName,
		GroupImage:     gc.GroupImage,
		LatestMessage:  (*database.LatestMessage)(gc.LatestMessage),
	}
}
