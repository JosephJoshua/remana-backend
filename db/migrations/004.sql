-- +migrate Up
ALTER TABLE login_codes
  DROP CONSTRAINT login_codes_user_id_fkey,
  ADD CONSTRAINT login_codes_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

ALTER TABLE roles
  DROP CONSTRAINT roles_store_id_fkey,
  ADD CONSTRAINT roles_store_id_fkey
    FOREIGN KEY (store_id) REFERENCES stores (store_id) ON DELETE CASCADE;

ALTER TABLE users
  DROP CONSTRAINT users_store_id_fkey,
  ADD CONSTRAINT users_store_id_fkey
    FOREIGN KEY (store_id) REFERENCES stores (store_id) ON DELETE CASCADE;

-- +migrate Down
ALTER TABLE login_codes
  DROP CONSTRAINT login_codes_user_id_fkey,
  ADD CONSTRAINT login_codes_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users (user_id);

ALTER TABLE roles
  DROP CONSTRAINT roles_store_id_fkey,
  ADD CONSTRAINT roles_store_id_fkey
    FOREIGN KEY (store_id) REFERENCES stores (store_id);

ALTER TABLE users
  DROP CONSTRAINT users_store_id_fkey,
  ADD CONSTRAINT users_store_id_fkey
    FOREIGN KEY (store_id) REFERENCES stores (store_id);
