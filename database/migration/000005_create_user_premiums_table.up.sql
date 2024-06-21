CREATE TABLE user_premiums (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_premiums_user_id ON user_premiums(user_id);
CREATE INDEX idx_user_premiums_start_at ON user_premiums(start_at);
CREATE INDEX idx_user_premiums_end_at ON user_premiums(end_at);
