-- name: CreateRole :exec
INSERT INTO roles (
  role_id,
  store_id,
  role_name,
  is_store_admin
)
VALUES (
  $1,
  $2,
  $3,
  $4
);

-- name: IsRoleNameTaken :one
SELECT 1
FROM roles
WHERE roles.store_id = $1 AND LOWER(roles.role_name) = LOWER(sqlc.arg('role_name'));
