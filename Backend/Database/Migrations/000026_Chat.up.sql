PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Chat (
    Id VARCHAR(36) NOT NULL,
    SenderId VARCHAR(36) NOT NULL,
    CreationDate VARCHAR(20) NOT NULL,
    Message TEXT NOT NULL,
    Image TEXT,
    ReceiverId VARCHAR(36) DEFAULT '',
    GroupId VARCHAR(36) DEFAULT '',

    PRIMARY KEY (Id),

    FOREIGN KEY (SenderId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
    FOREIGN KEY (ReceiverId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
    FOREIGN KEY (GroupId) REFERENCES "Groups"("Id") ON DELETE CASCADE,

    CHECK (
        (ReceiverId <> '' AND GroupId = '') OR
        (ReceiverId = '' AND GroupId <> '')
    )
);
