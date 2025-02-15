DROP INDEX IF EXISTS idx_transactions_sender_id;
DROP INDEX IF EXISTS idx_transactions_receiver_id;

DROP INDEX IF EXISTS idx_wallets_id;
DROP INDEX IF EXISTS idx_wallets_user_id;

ALTER TABLE transactions DROP CONSTRAINT fk_receiver;
ALTER TABLE transactions DROP CONSTRAINT fk_sender;

ALTER TABLE wallets DROP CONSTRAINT fk_wallet_user;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;

