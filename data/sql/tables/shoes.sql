CREATE TABLE IF NOT EXISTS shoes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brand TEXT NOT NULL, 
    name TEXT NOT NULL,
    silhouette TEXT, 
    image_url TEXT, 
    tags TEXT
);
