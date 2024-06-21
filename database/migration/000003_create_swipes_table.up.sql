CREATE TABLE swipes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    swipe_user_id INTEGER NOT NULL,
    swipe_like BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_swipes_user_id ON swipes(user_id);

CREATE INDEX idx_swipes_swipe_user_id ON swipes(swipe_user_id);
