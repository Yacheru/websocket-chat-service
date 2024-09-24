CREATE TABLE IF NOT EXISTS messages (
    id BIGINT PRIMARY KEY,
    player VARCHAR,
    player_uuid UUID,
    message TEXT,
    sent_at TIMESTAMP,
);
