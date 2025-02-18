ALTER TABLE inventory DROP CONSTRAINT fk_purchase_merch;
ALTER TABLE inventory DROP CONSTRAINT fk_purchase_user;

DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS merch;