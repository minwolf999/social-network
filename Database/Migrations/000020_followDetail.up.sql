PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS FollowDetail AS
  SELECT 
    f.Id,
    f.FollowerId AS FollowerId,
    User.Username AS Follower_Username,
    
    CASE
      WHEN CONCAT(User.FirstName, ' ', User.LastName) = ' ' THEN ''
      ELSE CONCAT(User.FirstName, ' ', User.LastName)
    END AS Follower_Name,

    User.ProfilePicture AS Follower_Picture,

    f.FollowedId AS FollowedId,
    Follower.Username AS Followed_Username,

    CASE
      WHEN CONCAT(Follower.FirstName, ' ', Follower.LastName) = ' ' THEN ''
      ELSE CONCAT(Follower.FirstName, ' ', Follower.LastName)
    END AS Followed_Name,

    Follower.ProfilePicture AS Followed_Picture

FROM Follower AS f
INNER JOIN UserInfo AS User ON User.Id = f.FollowerId
INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowedId
