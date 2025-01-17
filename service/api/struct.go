package api

import (
	"github.com/gofrs/uuid"
	"github.com/shernille37/WASAText/service/database"
)

type User struct {
	UserID   uuid.UUID `json:"userID"`
	Username string    `json:"username"`
	Image    *string   `json:"image"`
}

type PrivateConversation struct {
	User *User `json:"user"`
}

type GroupConversation struct {
	GroupName  string  `json:"groupName"`
	GroupImage *string `json:"groupImage"`
}

type LatestMessage struct {
	Message   string  `json:"message"`
	Image     *string `json:"image"`
	Timestamp string  `json:"timestamp"`
}

type Conversation struct {
	ConversationID uuid.UUID            `json:"conversationID"`
	Type           string               `json:"type"`
	Private        *PrivateConversation `json:"private"`
	Group          *GroupConversation   `json:"group"`
	LatestMessage  *LatestMessage       `json:"latestMessage"`
	Members        []User               `json:"members"`
}

type Message struct {
	MessageID          uuid.UUID  `json:"messageID"`
	SenderID           uuid.UUID  `json:"senderID"`
	SenderName         string     `json:"senderName"`
	ConversationID     uuid.UUID  `json:"conversationID"`
	ReplyMessageID     *uuid.UUID `json:"replyMessageID"`
	ReplyMessage       *string    `json:"replyMessage"`
	ReplyRecipientName *string    `json:"replyRecipientName"`
	Timestamp          string     `json:"timestamp"`
	HasImage           bool       `json:"hasImage"`
	MessageType        string     `json:"messageType"`
	MessageStatus      string     `json:"messageStatus"`
	TimeDelivered      *string    `json:"timeDelivered"`
	TimeRead           *string    `json:"timeRead"`
	Message            string     `json:"message"`
	Image              *string    `json:"image"`
	Reactions          []Reaction `json:"reactions"`
}

type Reader struct {
	User      *User  `json:"user"`
	Timestamp string `json:"timestamp"`
}

type Reaction struct {
	ReactionID *uuid.UUID `json:"reactionID"`
	Unicode    string     `json:"unicode"`
	Reactor    *User      `json:"reactor"`
	Count      *int       `json:"count"`
}

// -------- Mapping Functions --------

func (m *Message) ToDatabase() database.Message {

	var reactions []database.Reaction

	for _, reaction := range m.Reactions {
		reactions = append(reactions, database.Reaction{
			ReactionID: reaction.ReactionID,
			Unicode:    reaction.Unicode,
			Reactor:    (*database.User)(reaction.Reactor),
			Count:      reaction.Count,
		})
	}

	return database.Message{
		MessageID:          m.MessageID,
		SenderID:           m.SenderID,
		SenderName:         m.SenderName,
		ConversationID:     m.ConversationID,
		ReplyMessageID:     m.ReplyMessageID,
		ReplyMessage:       m.ReplyMessage,
		ReplyRecipientName: m.ReplyRecipientName,
		Timestamp:          m.Timestamp,
		HasImage:           m.HasImage,
		MessageType:        m.MessageType,
		MessageStatus:      m.MessageStatus,
		TimeDelivered:      m.TimeDelivered,
		Message:            m.Message,
		Image:              m.Image,
		Reactions:          reactions,
	}
}

func (m *Message) FromDatabase(mess database.Message) {

	var reactions = make([]Reaction, len(mess.Reactions))

	for idx := range reactions {
		reactions[idx] = Reaction{
			ReactionID: mess.Reactions[idx].ReactionID,
			Unicode:    mess.Reactions[idx].Unicode,
			Reactor:    (*User)(mess.Reactions[idx].Reactor),
			Count:      mess.Reactions[idx].Count,
		}
	}

	m.MessageID = mess.MessageID
	m.SenderID = mess.SenderID
	m.SenderName = mess.SenderName
	m.ConversationID = mess.ConversationID
	m.ReplyMessageID = mess.ReplyMessageID
	m.ReplyMessage = mess.ReplyMessage
	m.ReplyRecipientName = mess.ReplyRecipientName
	m.Timestamp = mess.Timestamp
	m.HasImage = mess.HasImage
	m.MessageType = mess.MessageType
	m.MessageStatus = mess.MessageStatus
	m.TimeDelivered = mess.TimeDelivered
	m.Message = mess.Message
	m.Image = mess.Image
	m.Reactions = reactions

}

func (u *User) FromDatabase(user database.User) {
	u.UserID = user.UserID
	u.Username = user.Username
	u.Image = user.Image
}

func (u *User) ToDatabase() database.User {
	return database.User{
		UserID:   u.UserID,
		Username: u.Username,
		Image:    u.Image,
	}
}

func (c *Conversation) ToDatabase() database.Conversation {
	var private *database.PrivateConversation
	var group *database.GroupConversation

	if c.Private != nil {

		private = &database.PrivateConversation{

			User: (*database.User)(c.Private.User),
		}
	}

	if c.Group != nil {

		group = &database.GroupConversation{

			GroupName:  c.Group.GroupName,
			GroupImage: c.Group.GroupImage,
		}
	}

	var members []database.User

	for _, member := range c.Members {
		members = append(members, database.User{
			UserID:   member.UserID,
			Username: member.Username,
			Image:    member.Image,
		})
	}

	return database.Conversation{
		ConversationID: c.ConversationID,
		Type:           c.Type,
		Private:        private,
		Group:          group,
		LatestMessage:  (*database.LatestMessage)(c.LatestMessage),
		Members:        members,
	}
}

func (c *Conversation) FromDatabase(conv database.Conversation) {
	c.ConversationID = conv.ConversationID
	c.Type = conv.Type
	c.LatestMessage = (*LatestMessage)(conv.LatestMessage)
	if conv.Private != nil {
		c.Private = &PrivateConversation{
			User: (*User)(conv.Private.User),
		}
	}

	if conv.Group != nil {
		c.Group = &GroupConversation{
			GroupName:  conv.Group.GroupName,
			GroupImage: conv.Group.GroupImage,
		}

	}

	var members = make([]User, len(conv.Members))
	for idx := range members {
		members[idx] = User{
			UserID:   conv.Members[idx].UserID,
			Username: conv.Members[idx].Username,
			Image:    conv.Members[idx].Image,
		}
	}

	c.Members = members

}

func (pc *PrivateConversation) FromDatabase(priv database.PrivateConversation) {

	pc.User = (*User)(priv.User)

}

func (pc *PrivateConversation) ToDatabase() database.PrivateConversation {
	return database.PrivateConversation{

		User: (*database.User)(pc.User),
	}
}

func (gc *GroupConversation) FromDatabase(group database.GroupConversation) {

	gc.GroupName = group.GroupName
	gc.GroupImage = group.GroupImage

}

func (gc *GroupConversation) ToDatabase() database.GroupConversation {
	return database.GroupConversation{

		GroupName:  gc.GroupName,
		GroupImage: gc.GroupImage,
	}
}

func (r *Reader) ToDatabase() database.Reader {

	return database.Reader{
		User: &database.User{
			UserID:   r.User.UserID,
			Username: r.User.Username,
			Image:    r.User.Image,
		},
		Timestamp: r.Timestamp,
	}
}

func (r *Reader) FromDatabase(rdb database.Reader) {
	r.User = (*User)(rdb.User)
	r.Timestamp = rdb.Timestamp
}

func (r *Reaction) ToDatabase() database.Reaction {

	return database.Reaction{
		ReactionID: r.ReactionID,
		Unicode:    r.Unicode,
		Reactor:    (*database.User)(r.Reactor),
	}
}

func (r *Reaction) FromDatabase(dr database.Reaction) {
	r.ReactionID = dr.ReactionID
	r.Unicode = dr.Unicode
	r.Reactor = (*User)(dr.Reactor)
}
