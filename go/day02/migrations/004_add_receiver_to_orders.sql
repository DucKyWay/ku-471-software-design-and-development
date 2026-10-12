ALTER TABLE orders ADD COLUMN receiver_name TEXT;
ALTER TABLE orders ADD COLUMN receiver_phone TEXT;

CREATE INDEX idx_orders_receiver_phone ON orders(receiver_phone);
