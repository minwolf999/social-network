PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS FollowDetail AS
  SELECT 
    f.Id,
    f.UserId AS UserId,
    User.Username AS User_Username,
    
    CASE
      WHEN CONCAT(User.FirstName, ' ', User.LastName) = ' ' THEN ''
      ELSE CONCAT(User.FirstName, ' ', User.LastName)
    END AS User_Name,

    f.FollowerId AS FollowerId,
    Follower.Username AS Follower_Username,

    CASE
      WHEN CONCAT(Follower.FirstName, ' ', Follower.LastName) = ' ' THEN ''
      ELSE CONCAT(Follower.FirstName, ' ', Follower.LastName)
    END AS Follower_Name

FROM Follower AS f
INNER JOIN UserInfo AS User ON User.Id = f.UserId
INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowerId
