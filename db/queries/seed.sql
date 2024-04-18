-- name: SeedStore :one
INSERT INTO stores (store_id, store_name, store_code, store_address, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING store_id;

-- name: SeedRole :one
INSERT INTO roles (role_id, role_name, store_id, is_store_admin)
VALUES ($1, $2, $3, $4)
RETURNING role_id;

-- name: SeedUser :one
INSERT INTO users (user_id, username, user_password, role_id, store_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id;

-- name: SeedLoginCode :one
INSERT INTO login_codes (login_code_id, user_id, login_code)
VALUES ($1, $2, $3)
RETURNING login_code_id;
