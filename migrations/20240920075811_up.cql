CREATE TABLE IF NOT EXISTS messages (
    player VARCHAR,
    player_uuid UUID,
    message TEXT,
    sent_at TIMESTAMP,
    PRIMARY KEY ( player, sent_at )
);
