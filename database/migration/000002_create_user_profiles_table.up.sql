CREATE TABLE user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    gender VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    picture_url VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_user_profiles_gender ON user_profiles(gender);
