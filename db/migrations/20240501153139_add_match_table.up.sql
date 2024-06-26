CREATE TABLE IF NOT EXISTS matchs (
    id SERIAL PRIMARY KEY,
    match_cat_id INT NOT NULL,
    user_cat_id INT NOT NULL,
    issued_by INT NOT NULL,
    accepted_by INT NOT NULL,
    is_approved BOOLEAN DEFAULT NULL,
    message TEXT NOT NULL,
    is_matched BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (issued_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (match_cat_id) REFERENCES cats(id) ON DELETE CASCADE,
    FOREIGN KEY (user_cat_id) REFERENCES cats(id) ON DELETE CASCADE
);