PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS FollowRequestDetail AS
    SELECT
        f.FollowerId AS FollowerId,
        CASE 
            WHEN User.Username = '' THEN CONCAT(User.FirstName, ' ', User.LastName)
            ELSE User.Username 
        END AS Follower_Name,

        f.FollowedId AS FollowedId,
        CASE 
            WHEN Follower.Username = '' THEN CONCAT(Follower.FirstName, ' ', Follower.LastName)
            ELSE Follower.Username 
        END AS Followed_Name
        

    FROM FollowingRequest AS f
    INNER JOIN UserInfo AS User ON User.Id = f.FollowerId
    INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowedId