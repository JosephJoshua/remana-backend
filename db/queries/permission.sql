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

-- name: DoesRoleExist :one
SELECT 1
FROM roles
WHERE roles.role_id = $1;

-- name: AssignPermissionsToRole :copyfrom
INSERT INTO role_permissions (
  role_id,
  permission_id
) VALUES (
  $1,
  $2
);
