CREATE TABLE IF NOT EXISTS global_config (
    id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at TEXT NOT NULL,
    scope TEXT NOT NULL,
    CONSTRAINT global_config_key_scope_key UNIQUE (key, scope)
);

CREATE TABLE IF NOT EXISTS Users (
    username VARCHAR(255) PRIMARY KEY,
    password_hash TEXT NOT NULL,
    campaign_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Campaigns table for storing campaign information
CREATE TABLE IF NOT EXISTS Campaigns (
    id VARCHAR(255) PRIMARY KEY,
    slot INT NOT NULL,
    valid_from TIMESTAMP NOT NULL,
    valid_to TIMESTAMP NOT NULL,    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Subscriptions table for storing subscription details
CREATE TABLE IF NOT EXISTS Subscriptions (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- UserCampaigns table for mapping users to campaigns
CREATE TABLE IF NOT EXISTS UserCampaigns (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL,
    campaign_id VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- UserSubscriptions table for tracking user subscriptions
CREATE TABLE IF NOT EXISTS UserSubscriptions (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    subscription_id VARCHAR(255) NOT NULL,
    campaign_id VARCHAR(255) NULL,
    status VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



INSERT INTO Subscriptions (id, name, price)
VALUES
    ('silver', 'Silver Plan', '100'),
    ('gold', 'Gold Plan', '200');