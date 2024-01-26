CREATE TABLE IF NOT EXISTS log
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level      TEXT not null,
    message    TEXT not null,
    created_at TEXT not null
);