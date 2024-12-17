PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS GroupDetail AS
  SELECT 
    g.Id,
    g.LeaderId,
    
    CASE 
      WHEN u.Username = '' THEN CONCAT(u.FirstName, ' ', u.LastName)
      ELSE u.Username 
    END AS Leader,

    g.MemberIds,
    g.groupName,
    g.GroupDescription,
    g.CreationDate,
    g.GroupPicture,
    g.Banner

FROM Groups AS g
INNER JOIN UserInfo AS u ON u.Id = g.LeaderId
