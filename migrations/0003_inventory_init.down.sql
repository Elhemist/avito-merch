DROP INDEX IF EXISTS idx_inventory_user_id_merch_id;
DROP INDEX IF EXISTS idx_inventory_user_id;

DROP INDEX IF EXISTS idx_merch_name;
DROP INDEX IF EXISTS idx_merch_id;

ALTER TABLE inventory DROP CONSTRAINT fk_purchase_merch;
ALTER TABLE inventory DROP CONSTRAINT fk_purchase_user;

DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS merch;