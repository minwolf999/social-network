CREATE TABLE IF NOT EXISTS Groups (
	Id VARCHAR(36) NOT NULL,
	LeaderId VARCHAR(36) NOT NULL,
	MemberIds TEXT NOT NULL,
	GroupName VARCHAR(200) NOT NULL,
	CreationDate VARCHAR(20) NOT NULL,

	PRIMARY KEY (Id),

	CONSTRAINT fk_leaderid FOREIGN KEY (LeaderId) REFERENCES "UserInfo"("Id")	
);