package database

const create_table = `

CREATE TABLE IF NOT EXISTS User(
	id INTEGER,
	username TEXT NOT NULL,
	image TEXT,

	PRIMARY KEY(id)
);


CREATE TABLE IF NOT EXISTS Conversation (
	id INTEGER,
	conversationType TEXT NOT NULL CHECK(conversationType IN ('personal', 'group')),
	groupName TEXT,
	groupImage TEXT,

	PRIMARY KEY(id),
	FOREIGN KEY (userID) REFERENCES User(id)
);

CREATE TABLE IF NOT EXISTS Members (
	userID INTEGER,
	conversationID INTEGER,

	PRIMARY KEY (userID, conversationID),
	FOREIGN KEY (userID) REFERENCES User(id),
	FOREIGN KEY (conversationID) REFERENCES Conversation(id)
);

CREATE TABLE IF NOT EXISTS Message (
	id INTEGER,
	senderID INTEGER NOT NULL,
	conversationID INTEGER NOT NULL,
	timestamp TEXT NOT NULL,
	messageType TEXT NOT NULL CHECK(messageType IN ('text', 'image')),
	messageStatus TEXT NOT NULL CHECK(messageStatus IN ('delivered', 'read')),
	timeDelivered TEXT,
	message TEXT,
	image TEXT,

	PRIMARY KEY (id),
	FOREIGN KEY (senderID) REFERENCES User(id),
	FOREIGN KEY (conversationID) REFERENCES Conversation(id)

);

CREATE TABLE IF NOT EXISTS Reader (
	userID INTEGER NOT NULL,
	messageID INTEGER NOT NULL,
	timestamp TEXT NOT NULL,

	PRIMARY KEY (userID, messageID)
);

CREATE TABLE IF NOT EXISTS Emoji (
	unicode TEXT NOT NULL,

	PRIMARY KEY (unicode)
);

CREATE TABLE IF NOT EXISTS Reaction (
	id INTEGER NOT NULL,
	userID INTEGER NOT NULL,
	messageID INTEGER NOT NULL,
	emoji TEXT NOT NULL,

	PRIMARY KEY (id),
	UNIQUE(userID, messageID, emoji),
	FOREIGN KEY (userID) REFERENCES User(id),
	FOREIGN KEY (messageID) REFERENCES Message(id),
	FOREIGN KEY (emoji) REFERENCES Emoji(unicode)
	
);

`

const database_structure = create_table