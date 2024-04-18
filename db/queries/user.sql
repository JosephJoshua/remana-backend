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
