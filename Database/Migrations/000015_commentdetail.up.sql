PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS CommentDetail AS
  SELECT 
    c.Id,
	c.Text,
	c.Image,
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
