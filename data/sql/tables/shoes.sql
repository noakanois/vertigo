CREATE TABLE IF NOT EXISTS shoes (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT,
    Subtitle TEXT,
    LastSale TEXT,
    ProductName TEXT UNIQUE,
    MainPicture TEXT,
    Attributes TEXT,
    Description TEXT,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    SpinningGifURL TEXT
);