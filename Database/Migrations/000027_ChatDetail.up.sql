PRAGMA foreign_keys = ON;

CREATE VIEW IF NOT EXISTS ChatDetail AS 
SELECT
    c.Id,
    c.SenderId,
    CASE 
        WHEN Sender.Username = '' THEN Sender.FirstName || ' ' || Sender.LastName
        ELSE Sender.Username 
    END AS Sender_Name,
    c.CreationDate,
    c.Message,
    c.Image,
    c.ReceiverId,
    CASE 
        WHEN c.ReceiverId <> '' THEN 
            CASE 
                WHEN Receiver.Username = '' THEN Receiver.FirstName || ' ' || Receiver.LastName
                ELSE Receiver.Username 
            END
        ELSE NULL
    END AS Receiver_Name,
    c.GroupId,
    CASE 
        WHEN c.GroupId <> '' THEN g.GroupName
        ELSE NULL
    END AS Group_Name
FROM Chat AS c
INNER JOIN UserInfo AS Sender ON Sender.Id = c.SenderId
LEFT JOIN UserInfo AS Receiver ON Receiver.Id = c.ReceiverId
LEFT JOIN Groups AS g ON g.Id = c.GroupId;
