PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Groups (
	Id VARCHAR(36) NOT NULL,
	LeaderId VARCHAR(36) NOT NULL,
	MemberIds TEXT NOT NULL,
	GroupName VARCHAR(200) NOT NULL,
	GroupDescription VARCHAR(500),
	CreationDate VARCHAR(20) NOT NULL,
	Banner TEXT,
	GroupPicture TEXT,

	PRIMARY KEY (Id),

	CONSTRAINT fk_leaderid FOREIGN KEY (LeaderId) REFERENCES "UserInfo"("Id")	
);