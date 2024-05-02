CREATE TABLE IF NOT EXISTS cats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    race VARCHAR(100),
    sex VARCHAR(100),
    age_in_months INT,
    description TEXT,
    image_urls TEXT[],
    user_id INT NOT NULL,
    has_match BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);