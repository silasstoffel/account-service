CREATE TABLE IF NOT EXISTS accounts(
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(30),
    last_name VARCHAR(50),
    full_name VARCHAR(80),
    email VARCHAR(60) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT true,
    hashed_pwd VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_account_email ON accounts(email);
