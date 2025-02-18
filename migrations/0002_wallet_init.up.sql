CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    balance INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE wallets ADD CONSTRAINT fk_wallet_user FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE transactions ADD CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users(id);
ALTER TABLE transactions ADD CONSTRAINT fk_receiver FOREIGN KEY (receiver_id) REFERENCES users(id);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);

CREATE INDEX idx_transactions_receiver_id ON transactions(receiver_id);
CREATE INDEX idx_transactions_sender_id ON transactions(sender_id);