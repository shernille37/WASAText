package database

type PrivateConversation struct {
	ConversationID uint64
	UserID string
	Name string
	Image string
	LatestMessage LatestMessage

}

type GroupConversation struct {
	ConversationID uint64
	GroupName string
	GroupImage string
	LatestMessage *LatestMessage
}


type LatestMessage struct {
	MessageType bool
	Message string
	Timestamp string
}

type Conversation struct {
	Type string
	Private *PrivateConversation
	GroupConversation *GroupConversation
}
