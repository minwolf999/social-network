PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Chat (
	Id VARCHAR(36) NOT NULL,
	SenderId VARCHAR(36) NOT NULL,
	CreationDate VARCHAR(20) NOT NULL,
	Message TEXT NOT NULL,
	Image TEXT,
	ReceiverId VARCHAR(36),
	GroupId VARCHAR(36),

	PRIMARY KEY (Id),

	CONSTRAINT fk_senderid FOREIGN KEY (SenderId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
)