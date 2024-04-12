-- +migrate Up
CREATE TABLE stores (
  store_id UUID NOT NULL PRIMARY KEY,
  store_name TEXT NOT NULL,
  store_code TEXT NOT NULL,
  store_address TEXT NOT NULL,
  phone_number TEXT NOT NULL
);

CREATE TABLE roles (
  role_id UUID NOT NULL PRIMARY KEY,
  role_name TEXT NOT NULL,
  store_id UUID NOT NULL REFERENCES stores (store_id),
  is_store_admin BOOLEAN NOT NULL
);

CREATE TABLE users (
  user_id UUID NOT NULL PRIMARY KEY,
  username TEXT NOT NULL,
  user_password TEXT NOT NULL,
  role_id UUID NOT NULL REFERENCES roles (role_id),
  store_id UUID NOT NULL REFERENCES stores (store_id)
);

-- +migrate Down
DROP TABLE stores;
DROP TABLE roles;
DROP TABLE users;
