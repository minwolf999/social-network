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
