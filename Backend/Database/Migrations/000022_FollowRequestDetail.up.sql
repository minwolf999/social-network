PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS FollowRequestDetail AS
    SELECT
        f.UserId AS UserId,
        CASE 
            WHEN User.Username = '' THEN CONCAT(User.FirstName, ' ', User.LastName)
            ELSE User.Username 
        END AS User_Name,

        f.FollowerId AS FollowerId,
        CASE 
            WHEN Follower.Username = '' THEN CONCAT(Follower.FirstName, ' ', Follower.LastName)
            ELSE Follower.Username 
        END AS Follower_Name
        

    FROM FollowingRequest AS f
    INNER JOIN UserInfo AS User ON User.Id = f.UserId
    INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowerId