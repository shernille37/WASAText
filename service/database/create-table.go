package database

const create_table = `

CREATE TABLE IF NOT EXISTS User(
	pk INTEGER,
	userID TEXT NOT NULL,
	username TEXT NOT NULL,
	image TEXT,

	PRIMARY KEY(pk),
	UNIQUE(userID),
	UNIQUE(username)
);


CREATE TABLE IF NOT EXISTS Conversation (
	pk INTEGER,
	conversationID TEXT NOT NULL,
	conversationType TEXT NOT NULL CHECK(conversationType IN ('private', 'group')),
	groupName TEXT,
	groupImage TEXT,

	PRIMARY KEY(pk),
	UNIQUE(conversationID)
);

CREATE TABLE IF NOT EXISTS Members (
	userID TEXT NOT NULL,
	conversationID TEXT NOT NULL,

	PRIMARY KEY (userID, conversationID),
	FOREIGN KEY (userID) REFERENCES User(userID) ON DELETE CASCADE,
	FOREIGN KEY (conversationID) REFERENCES Conversation(conversationID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Message (
	pk INTEGER,
	messageID TEXT NOT NULL,
	senderID TEXT NOT NULL,
	conversationID TEXT NOT NULL,
	replyMessageID TEXT,
	timestamp TEXT NOT NULL DEFAULT current_timestamp,
	messageType TEXT NOT NULL CHECK(messageType IN ('default', 'reply', 'forward')) DEFAULT 'default',
	hasImage INTEGER DEFAULT 0,
	messageStatus TEXT NOT NULL CHECK(messageStatus IN ('delivered', 'read', 'sent')) DEFAULT 'sent',
	timeRead TEXT,
	timeDelivered TEXT,
	message TEXT NOT NULL,
	image TEXT,

	PRIMARY KEY (pk),
	UNIQUE (messageID),
	FOREIGN KEY (replyMessageID) REFERENCES Message(messageID) ON DELETE SET NULL,
	FOREIGN KEY (senderID) REFERENCES User(userID) ON DELETE CASCADE,
	FOREIGN KEY (conversationID) REFERENCES Conversation(conversationID) ON DELETE CASCADE

);

CREATE TABLE IF NOT EXISTS Reader (
	userID TEXT NOT NULL,
	messageID TEXT NOT NULL,
	timestamp TEXT NOT NULL,

	PRIMARY KEY (userID, messageID)
);

CREATE TABLE IF NOT EXISTS Emoji (
	unicode TEXT NOT NULL,

	PRIMARY KEY (unicode)
);

CREATE TABLE IF NOT EXISTS Reaction (
	pk INTEGER NOT NULL,
	reactionID TEXT NOT NULL,
	userID TEXT NOT NULL,
	messageID TEXT NOT NULL,
	emoji TEXT NOT NULL,

	PRIMARY KEY (pk),
	UNIQUE(reactionID),
	UNIQUE(userID, messageID),
	FOREIGN KEY (userID) REFERENCES User(userID) ON DELETE CASCADE,
	FOREIGN KEY (messageID) REFERENCES Message(messageID) ON DELETE CASCADE,
	FOREIGN KEY (emoji) REFERENCES Emoji(unicode) ON DELETE CASCADE
	
);

`
