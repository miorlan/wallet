CREATE TABLE IF NOT EXISTS wallets (
                         wallet_id UUID PRIMARY KEY,
                         balance DECIMAL(15, 2) NOT NULL DEFAULT 0.0
);