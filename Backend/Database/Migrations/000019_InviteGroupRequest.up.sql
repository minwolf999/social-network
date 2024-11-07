PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS InviteGroupRequest (
	SenderId VARCHAR(36) NOT NULL,
	GroupId VARCHAR(36) NOT NULL,
	ReceiverId VARCHAR(36) NOT NULL,

	CONSTRAINT fk_senderid FOREIGN KEY (SenderId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
	CONSTRAINT fk_followerid FOREIGN KEY (GroupId) REFERENCES "Groups"("Id") ON DELETE CASCADE,
	CONSTRAINT fk_receiverid FOREIGN KEY (ReceiverId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
)