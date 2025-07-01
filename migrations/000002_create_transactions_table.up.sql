CREATE TABLE IF NOT EXISTS transactions (
                              id SERIAL PRIMARY KEY,
                              wallet_id UUID REFERENCES wallets(wallet_id),
                              operation_type VARCHAR(10) NOT NULL,
                              amount DECIMAL(15, 2) NOT NULL,
                              created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

);