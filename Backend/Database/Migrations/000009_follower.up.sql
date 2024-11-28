PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Follower (
	Id VARCHAR(36) NOT NULL,
	UserId VARCHAR(36) NOT NULL,
	FollowedId VARCHAR(36) NOT NULL,

	PRIMARY KEY (Id),

	CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
	CONSTRAINT fk_followedid FOREIGN KEY (FollowedId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
);