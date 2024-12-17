PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS FollowingRequest (
	FollowerId VARCHAR(36) NOT NULL,
	FollowedId VARCHAR(36) NOT NULL,

	CONSTRAINT fk_followerid FOREIGN KEY (FollowerId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
	CONSTRAINT fk_followedid FOREIGN KEY (FollowedId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
)