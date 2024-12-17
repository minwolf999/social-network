PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Event (
	Id VARCHAR(36),
	GroupId VARCHAR(36),
	OrganisatorId VARCHAR(36),
	Title VARCHAR(200),
	Description VARCHAR(1000),
	DateOfTheEvent VARCHAR(20),
	Image TEXT,

	PRIMARY KEY (Id),

	CONSTRAINT fk_groupid FOREIGN KEY (GroupId) REFERENCES "Groups"("Id") ON DELETE CASCADE,
	CONSTRAINT fk_organisatorid FOREIGN KEY (OrganisatorId) REFERENCES "UserInfo"("Id")
);