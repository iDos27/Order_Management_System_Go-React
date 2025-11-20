CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Wstawienie przykładowych użytkowników
-- Hasła: admin123 i gosc123 (zahashowane bcrypt)
INSERT INTO users (email, password_hash, role) VALUES
('admin@test.com', '$2a$10$rVHJ3J7VXD2eN8fFZJ5pF.yxHf8Y0KvN0CKNdN5yN.8KJN5pF.yxH', 'admin'),
('gosc@test.com', '$2a$10$rVHJ3J7VXD2eN8fFZJ5pF.yxHf8Y0KvN0CKNdN5yN.8KJN5pF.yxH', 'employee')
ON CONFLICT (email) DO NOTHING;
