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

    GROUP_CONCAT(DISTINCT u2.Id) AS JoinUsers,

    GROUP_CONCAT(DISTINCT u3.Id) AS DeclineUsers

FROM Event AS e
INNER JOIN Groups AS g ON g.Id = e.GroupId
INNER JOIN UserInfo AS u1 ON u1.Id = e.OrganisatorId

LEFT JOIN JoinEvent AS j ON j.EventId = e.Id
LEFT JOIN UserInfo AS u2 ON u2.Id = j.UserId

LEFT JOIN DeclineEvent AS d ON d.EventId = e.Id
LEFT JOIN UserInfo AS u3 ON u3.Id = d.UserId

GROUP BY 
    e.Id, g.Id, g.GroupName, u1.Username, u1.FirstName, u1.LastName, 
    e.Title, e.Description, e.DateOfTheEvent;
