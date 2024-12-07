PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS EventDetail AS
  SELECT 
    e.Id,
    g.Id AS GroupId,
    g.GroupName,

    CASE 
      WHEN u1.Username = '' THEN CONCAT(u1.FirstName, ' ', u1.LastName)
      ELSE u1.Username 
    END AS Organisator,

    e.Title,
    e.Description,
    e.DateOfTheEvent,

    GROUP_CONCAT(DISTINCT CASE 
        WHEN u2.Username = '' THEN CONCAT(u2.FirstName, ' ', u2.LastName)
        ELSE u2.Username 
    END) AS JoinUsers,

    GROUP_CONCAT(DISTINCT CASE
        WHEN u3.Username = '' THEN CONCAT(u3.FirstName, ' ', u3.LastName)
        ELSE u3.Username
    END) AS DeclineUsers


FROM Event AS e
INNER JOIN Groups AS g ON g.Id = e.GroupId
INNER JOIN UserInfo AS u1 ON u1.Id = e.OrganisatorId

LEFT JOIN JoinEvent AS j ON j.EventId = e.Id
LEFT JOIN UserInfo AS u2 ON u2.Id = j.UserId

LEFT JOIN DeclineEvent AS d ON d.EventId = e.Id
LEFT JOIN UserInfo AS u3 ON u3.Id = d.UserId