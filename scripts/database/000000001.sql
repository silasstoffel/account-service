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


CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(40) PRIMARY KEY,
    data_id VARCHAR(40) NULL,
    occurred_at TIMESTAMP NOT NULL,
    "type" VARCHAR(80) NOT NULL,
    "source" VARCHAR(80) NOT NULL,
    "data" JSON NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_event_type ON events ("type");
CREATE INDEX IF NOT EXISTS idx_event_source ON events ("source");
CREATE INDEX IF NOT EXISTS idx_event_data_id ON events ("data_id");


CREATE TABLE IF NOT EXISTS webhook_subscriptions(
    id VARCHAR(40) PRIMARY KEY,
    event_type VARCHAR(80) NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_event_type ON webhook_subscriptions ("event_type");
