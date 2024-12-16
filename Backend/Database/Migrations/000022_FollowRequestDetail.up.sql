PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS FollowRequestDetail AS
    SELECT
        f.FollowerId AS FollowerId,
        CASE 
            WHEN User.Username = '' THEN CONCAT(User.FirstName, ' ', User.LastName)
            ELSE User.Username 
        END AS Follower_Name,
        User.ProfilePicture AS Follower_Picture,

        f.FollowedId AS FollowedId,
        CASE 
            WHEN Follower.Username = '' THEN CONCAT(Follower.FirstName, ' ', Follower.LastName)
            ELSE Follower.Username 
        END AS Followed_Name,
        Follower.ProfilePicture AS Followed_Picture
        

    FROM FollowingRequest AS f
    INNER JOIN UserInfo AS User ON User.Id = f.FollowerId
    INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowedId