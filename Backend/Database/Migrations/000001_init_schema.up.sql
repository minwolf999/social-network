CREATE TABLE IF NOT EXISTS Auth (
	Id VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
	Email VARCHAR(100) NOT NULL UNIQUE,
	Password VARCHAR(50) NOT NULL,
	ConnectionAttempt INTEGER
);

CREATE TABLE IF NOT EXISTS UserInfo (
	Id VARCHAR(36) NOT NULL UNIQUE REFERENCES "Auth"("Id"),
	Email VARCHAR(100) NOT NULL UNIQUE REFERENCES "Auth"("Email"),
	FirstName VARCHAR(50) NOT NULL, 
	LastName VARCHAR(50) NOT NULL,
	BirthDate VARCHAR(20) NOT NULL,
	ProfilePicture VARCHAR(400000),
	Username VARCHAR(50),
	AboutMe VARCHAR(280)
);

CREATE TABLE IF NOT EXISTS Post (
	Id VARCHAR(36) NOT NULL UNIQUE,
	AuthorId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
	Text VARCHAR(1000) NOT NULL,
	Image VARCHAR(100),
	CreationDate VARCHAR(20) NOT NULL,
	IsGroup VARCHAR(36) REFERENCES "Groups"("Id"),
	LikeCount INTEGER,
	DislikeCount INTEGER
);

CREATE TABLE IF NOT EXISTS LikePost (
	PostId VARCHAR(36) NOT NULL REFERENCES "Post"("Id"),
	UserId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id")
);

CREATE TABLE IF NOT EXISTS DislikePost (
	PostId VARCHAR(36) NOT NULL REFERENCES "Post"("Id"),
	UserId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id")
);

CREATE TABLE IF NOT EXISTS Comment (
	Id VARCHAR(36) NOT NULL UNIQUE,
	AuthorId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
	Text VARCHAR(1000) NOT NULL,
	CreationDate VARCHAR(20) NOT NULL,

	PostId VARCHAR(36) REFERENCES "Post"("Id"),

	LikeCount INTEGER,
	DislikeCount INTEGER
);

CREATE TABLE IF NOT EXISTS LikeComment (
	PostId VARCHAR(36) NOT NULL REFERENCES "Comment"("Id"),
	UserId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id")
);

CREATE TABLE IF NOT EXISTS DislikeComment (
	PostId VARCHAR(36) NOT NULL REFERENCES "Comment"("Id"),
	UserId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id")
);

CREATE TABLE IF NOT EXISTS Follower (
	Id VARCHAR(36) NOT NULL UNIQUE,
	UserId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
	FollowerId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id")
);

CREATE TABLE IF NOT EXISTS Groups (
	Id VARCHAR(36) NOT NULL UNIQUE
);


CREATE VIEW PostDetail AS
  SELECT 
    p.Id,
	p.Text,
	p.Image,
	p.CreationDate,
	p.IsGroup,
	p.AuthorId,
	p.LikeCount,
	p.DislikeCount,
	u.FirstName,
	u.LastName,
	u.ProfilePicture,
	u.Username
FROM Post AS p
INNER JOIN UserInfo AS u ON p.AuthorId = u.Id;

CREATE VIEW CommentDetail AS
  SELECT 
    c.Id,
	c.Text,
	c.CreationDate,
	c.AuthorId,
	c.LikeCount,
	c.DislikeCount,
	c.PostId,
	u.FirstName,
	u.LastName,
	u.ProfilePicture,
	u.Username
FROM Comment AS c
INNER JOIN UserInfo AS u ON c.AuthorId = u.Id;