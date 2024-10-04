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