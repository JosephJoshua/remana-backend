-- +migrate Up
CREATE TABLE stores (
  store_id uuid NOT NULL PRIMARY KEY,
  store_name TEXT NOT NULL,
  store_code TEXT NOT NULL,
  store_address TEXT NOT NULL,
  phone_number TEXT NOT NULL
);

-- +migrate Down
DROP TABLE stores;
