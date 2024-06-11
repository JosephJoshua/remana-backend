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

-- name: HasPermission :one
SELECT COUNT(*)
FROM role_permissions
LEFT JOIN
  permissions ON permissions.permission_id = role_permissions.permission_id
LEFT JOIN
  permission_groups ON permission_groups.permission_group_id = permissions.permission_group_id
WHERE
  role_permissions.role_id = $1 AND
  permissions.permission_name = $2 AND
  permission_groups.permission_group_name = $3;

-- name: IsStoreAdmin :one
SELECT roles.is_store_admin
FROM roles
WHERE roles.role_id = $1;
