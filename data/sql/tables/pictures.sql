CREATE TABLE IF NOT EXISTS pictures (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    LocalLocation TEXT,
    DiscordImageLink TEXT,
    DiscordMessageId TEXT,
    Latitude REAL,
    Longitude REAL,
    TakenAt DATETIME,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);

