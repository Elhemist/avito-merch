ALTER TABLE transactions DROP CONSTRAINT fk_receiver;
ALTER TABLE transactions DROP CONSTRAINT fk_sender;

ALTER TABLE wallets DROP CONSTRAINT fk_wallet_user;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS wallets;

