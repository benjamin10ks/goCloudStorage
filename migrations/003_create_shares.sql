CREATE TABLE IF NOT EXISTS shares (
    id SERIAL PRIMARY KEY,
    file_id INT REFERENCES files (id) ON DELETE CASCADE,
    shared_by INTEGER REFERENCES users (id) ON DELETE CASCADE,
    shared_with INTEGER REFERENCES users (id) ON DELETE CASCADE,
    permission VARCHAR(10) DEFAULT 'read', -- 'read' or 'write'
    created_at TIMESTAMP DEFAULT NOW()
);
