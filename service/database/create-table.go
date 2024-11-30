package database

const create_table = `

CREATE TABLE IF NOT EXISTS User(
	pk INTEGER,
	userID TEXT NOT NULL,
	username TEXT NOT NULL,
	image TEXT,

	PRIMARY KEY(pk),
	UNIQUE(userID)
);


CREATE TABLE IF NOT EXISTS Conversation (
	pk INTEGER,
	conversationID TEXT NOT NULL,
	conversationType TEXT NOT NULL CHECK(conversationType IN ('personal', 'group')),
	groupName TEXT,
	groupImage TEXT,

	PRIMARY KEY(pk),
	UNIQUE(conversationID)
);

CREATE TABLE IF NOT EXISTS Members (
	userID TEXT NOT NULL,
	conversationID TEXT NOT NULL,

	PRIMARY KEY (userID, conversationID),
	FOREIGN KEY (userID) REFERENCES User(userID),
	FOREIGN KEY (conversationID) REFERENCES Conversation(conversationID)
);

CREATE TABLE IF NOT EXISTS Message (
	pk INTEGER,
	messageID TEXT NOT NULL,
	senderID TEXT NOT NULL,
	conversationID TEXT NOT NULL,
	timestamp TEXT NOT NULL,
	messageType TEXT NOT NULL CHECK(messageType IN ('text', 'image')),
	messageStatus TEXT NOT NULL CHECK(messageStatus IN ('delivered', 'read')),
	timeDelivered TEXT,
	message TEXT,
	image TEXT,

	PRIMARY KEY (pk),
	UNIQUE (messageID),
	FOREIGN KEY (senderID) REFERENCES User(userID),
	FOREIGN KEY (conversationID) REFERENCES Conversation(conversationID)

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
	UNIQUE(userID, messageID, emoji),
	FOREIGN KEY (userID) REFERENCES User(userID),
	FOREIGN KEY (messageID) REFERENCES Message(messageID),
	FOREIGN KEY (emoji) REFERENCES Emoji(unicode)
	
);

`

const database_structure = create_table