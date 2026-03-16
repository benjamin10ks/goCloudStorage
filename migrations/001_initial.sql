CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    storage_quota BIGINT DEFAULT 1073741824, -- Default 1GB
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    filepath VARCHAR(500) NOT NULL,
    filesize BIGINT NOT NULL,
    mime_type VARCHAR(100),
    uploaded_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()

)

CREATE TABLE IF NOT EXISTS shares (
    id SERIAL PRIMARY KEY,
    file_id INT REFERENCES files (id) ON DELETE CASCADE,
    shared_by INTEGER REFERENCES users (id) ON DELETE CASCADE,
    shared_with INTEGER REFERENCES users (id) ON DELETE CASCADE,
    permission VARCHAR(10) DEFAULT 'read', -- 'read' or 'write'
    created_at TIMESTAMP DEFAULT NOW()
);
