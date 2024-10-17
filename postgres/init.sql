-- init.sql
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    user VARCHAR(255) NOT NULL,
    item_id VARCHAR(255) NOT NULL
);
