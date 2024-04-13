-- name: GetUserByUsernameAndStoreCode :one
SELECT users.user_id, users.user_password, roles.is_store_admin
FROM users
LEFT JOIN stores ON stores.store_id = users.store_id
LEFT JOIN roles ON roles.role_id = users.role_id
WHERE users.username = $1 AND stores.store_code = $2
LIMIT 1;

-- name: GetLoginCodeByUserIDAndCode :one
SELECT login_codes.login_code_id
FROM login_codes
WHERE login_codes.user_id = $1 AND login_codes.login_code = $2;

-- name: DeleteLoginCodeByID :exec
DELETE FROM login_codes
WHERE login_codes.login_code_id = $1;

-- name: GetUserDetailsByID :one
SELECT
  users.user_id,
  users.username,
  roles.role_id,
  roles.role_name,
  roles.is_store_admin,
  stores.store_id,
  stores.store_name,
  stores.store_code
FROM users
LEFT JOIN stores ON stores.store_id = users.store_id
LEFT JOIN roles ON roles.role_id = users.role_id
WHERE users.user_id = $1
LIMIT 1;
