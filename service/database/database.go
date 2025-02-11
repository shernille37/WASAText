/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error
	Close() error

	Login(user User) (User, error)
	ListConversation(id uuid.UUID) ([]Conversation, error)
	GetConversation(id uuid.UUID, conversationID uuid.UUID) (Conversation, error)

	ListPrivateConversation(id uuid.UUID) ([]Conversation, error)
	AddPrivateChat(senderID uuid.UUID, mpb MessagePrivateBody) (Conversation, error)

	ListGroupConversation(id uuid.UUID) ([]Conversation, error)
	AddGroupConversation(senderID uuid.UUID, mgb MessageGroupBody) (Conversation, error)
	UpdateGroupName(conversationID uuid.UUID, newGroupName string) error
	UpdateGroupImage(conversationID uuid.UUID, newGroupPhoto string) error
	ListGroupMembers(conversationID uuid.UUID) ([]User, error)
	AddGroupMembers(conversationID uuid.UUID, gmb GroupMemberBody) ([]User, error)
	LeaveGroupConversation(userID uuid.UUID, conversationID uuid.UUID) error

	ListMessages(conversationID uuid.UUID) ([]Message, error)
	AddMessage(senderID uuid.UUID, conversationID uuid.UUID, mb MessageBody) (Message, error)
	ListReaders(conversationID uuid.UUID, messageID uuid.UUID) ([]Reader, error)
	ForwardMessage(senderID uuid.UUID, messageID uuid.UUID, fmb ForwardMessageBody) (Conversation, error)
	DeleteMessage(conversationID uuid.UUID, messageID uuid.UUID) error
	UpdateMessageToDelivered(senderID uuid.UUID, conversationID uuid.UUID) error
	UpdateMessageToRead(senderID uuid.UUID, conversationID uuid.UUID) error

	ListReactions(conversationID uuid.UUID, messageID uuid.UUID) ([]Reaction, error)
	AddReaction(userID uuid.UUID, messageID uuid.UUID, rb ReactionBody) (Reaction, error)
	DeleteReaction(reactionID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) error

	ListUsers(id uuid.UUID) ([]User, error)
	UpdateUsername(userID uuid.UUID, newUsername string) error
	UpdateUserImage(userID uuid.UUID, newUserImage string) error

	// Utils
	GetUserByID(id string) (User, error)
	CheckConversationMembership(id uuid.UUID, members []uuid.UUID) error
	CheckMessageMembership(conversationID uuid.UUID, messageID uuid.UUID) error
	CheckExistingConversation(senderID uuid.UUID, receiverID uuid.UUID) error
	CheckMessageOwnership(userID uuid.UUID, conversationID uuid.UUID, messageID uuid.UUID) error
	CheckReactionOwnership(userID uuid.UUID, reactionID uuid.UUID) error
	ListEmojis() ([]string, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// CREATE THE DATABASE STRUCTURE
	_, err := db.Exec(database_structure)
	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func (db *appdbimpl) Close() error {
	return db.c.Close()
}
