CREATE TABLE IF NOT EXISTS Comment (
	Id VARCHAR(36) NOT NULL,
	AuthorId VARCHAR(36) NOT NULL,
	Text VARCHAR(1000) NOT NULL,
	CreationDate VARCHAR(20) NOT NULL,
	PostId VARCHAR(36),
	LikeCount INTEGER,
	DislikeCount INTEGER,

	PRIMARY KEY (Id),

	CONSTRAINT fk_authorid FOREIGN KEY (AuthorId) REFERENCES "UserInfo"("Id"),
	CONSTRAINT fk_postid FOREIGN KEY (PostId) REFERENCES "Post"("Id")
);