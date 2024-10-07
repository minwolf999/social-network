PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Post (
    Id VARCHAR(36) NOT NULL,
    AuthorId VARCHAR(36) NOT NULL,
    Text VARCHAR(1000) NOT NULL,
    Image TEXT,
    CreationDate VARCHAR(20) NOT NULL,
    Status TEXT NOT NULL,
    IsGroup VARCHAR(36),
    LikeCount INTEGER,
    DislikeCount INTEGER,

	PRIMARY KEY (Id),

	CONSTRAINT fk_authorid FOREIGN KEY (AuthorId) REFERENCES "UserInfo"("Id"),
	CONSTRAINT fk_isgroup FOREIGN KEY (IsGroup) REFERENCES "Groups"("Id") ON DELETE CASCADE
);