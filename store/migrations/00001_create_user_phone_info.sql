-- +goose Up
CREATE TABLE IF NOT EXISTS USER_PHONE_INFO (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id_str TEXT,
    screen_name TEXT,
    blog_text_raw TEXT,
    blog_region_name TEXT,
    blog_source TEXT,
    blog_created_at TEXT,
    blog_id_str TEXT,
    blog_mblog_id TEXT,
    user_ip_location TEXT,
    user_created_at TEXT,
    gender TEXT,
    phone_type TEXT,
    phone_brand TEXT,
    batch_id TEXT
);

-- +goose Down
DROP TABLE USER_PHONE_INFO;
