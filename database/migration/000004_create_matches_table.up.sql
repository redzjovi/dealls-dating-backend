CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user_id_1 INTEGER NOT NULL,
    user_id_2 INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_matches_user_id_1 ON matches(user_id_1);

CREATE INDEX idx_matches_user_id_2 ON matches(user_id_2);

