PRAGMA foreign_keys = ON;

CREATE View IF NOT EXISTS JoinGroupRequestDetail AS
	SELECT
		j.UserId,
		
		CASE 
            WHEN u.Username = '' THEN CONCAT(u.FirstName, ' ', u.LastName)
            ELSE u.Username 
        END AS User_Name,

		j.GroupId,
		g.GroupName

	FROM JoinGroupRequest AS j
	INNER JOIN UserInfo AS u ON u.Id = j.UserId
	INNER JOIN Groups AS g ON g.Id = j.GroupId