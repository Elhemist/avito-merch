CREATE TABLE IF NOT EXISTS merch (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    merch_item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1
);

ALTER TABLE inventory ADD CONSTRAINT fk_purchase_user FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE inventory ADD CONSTRAINT fk_purchase_merch FOREIGN KEY (merch_item_id) REFERENCES merch(id);

CREATE INDEX idx_merch_name ON merch(name);

CREATE INDEX idx_inventory_user_id ON inventory(user_id);
CREATE INDEX idx_inventory_user_id_merch_id ON inventory(user_id, merch_item_id);

INSERT INTO merch (name, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);