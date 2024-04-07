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

-- Password: super.user
insert into accounts
    (id, "name", last_name, full_name, email, phone, hashed_pwd)
values
    ('01HTXANJXPMYPW89E6ZX0NFC63', 'Super', 'User', 'Super user', 'super.user@gmail.com', '+5511999999999', '$2a$15$ul/Uy2i4BhomUvDZFqPX4O7cXT3f06UEn3N9jmk3UsxTyCM.dSMfS');

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

CREATE TABLE IF NOT EXISTS webhook_transactions(
    id VARCHAR(36) PRIMARY KEY,
    event_id VARCHAR(36) NOT NULL,
    subscription_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(80),
    received_status_code smallint DEFAULT NULL,
    started_at TIMESTAMP NULL,
    finished_at TIMESTAMP DEFAULT NULL,
    number_of_requests smallint DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_webhook_transactions_event_type ON webhook_transactions("event_type");
CREATE INDEX IF NOT EXISTS idx_webhook_transactions_event_id ON webhook_transactions("event_id");


CREATE TABLE IF NOT EXISTS permissions(
    id VARCHAR(36) PRIMARY KEY,
    scope VARCHAR(80),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

insert into permissions (id, scope) values ('01HTXA4GGBDZY3S1E7SNSRQ6MG', 'account-service:*');
insert into permissions (id, scope) values ('01HTXA59W476S0MHK8MSD8ZKBC', 'account-service:create-account');
insert into permissions (id, scope) values ('01HTXA5GGK0WV9PDZVCN2C1N5K', 'account-service:update-account');
insert into permissions (id, scope) values ('01HTXA5QBA26MNW63435XS25PK', 'account-service:list-accounts');
insert into permissions (id, scope) values ('01HTXA5X7CRCSMGH1QPVDR86EM', 'account-service:get-account');
insert into permissions (id, scope) values ('01HTXA84W1SYGA35VVV5SYQZXA', 'account-service:disable-account');
insert into permissions (id, scope) values ('01HTXA90ZK1XZKWX2RGMB05CAJ', 'account-service:enable-account');


CREATE TABLE IF NOT EXISTS account_permissions (
    permission_id VARCHAR(36),
    account_id VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    PRIMARY KEY (account_id, permission_id)
);

insert into account_permissions (account_id, permission_id) values ('01HTXANJXPMYPW89E6ZX0NFC63', '01HTXA4GGBDZY3S1E7SNSRQ6MG');
