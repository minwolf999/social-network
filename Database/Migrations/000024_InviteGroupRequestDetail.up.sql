PRAGMA foreign_keys = ON;

CREATE View IF NOT EXISTS InviteGroupRequestDetail AS
	SELECT
		i.SenderId,
		
		CASE 
            WHEN Sender.Username = '' THEN CONCAT(Sender.FirstName, ' ', Sender.LastName)
            ELSE Sender.Username 
        END AS Sender_Name,

		i.GroupId,
		g.GroupName,

		i.ReceiverId,
		
		CASE
            WHEN Receiver.Username = '' THEN CONCAT(Receiver.FirstName, ' ', Receiver.LastName)
            ELSE Receiver.Username 
        END AS Receiver_Name
		
	FROM InviteGroupRequest AS i
	INNER JOIN UserInfo AS Sender ON Sender.Id = i.SenderId
	INNER JOIN Groups AS g ON g.Id = i.GroupId
	INNER JOIN UserInfo AS Receiver ON Receiver.Id = i.ReceiverId