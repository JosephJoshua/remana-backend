-- name: CreateDamageType :exec
INSERT INTO damage_types (
  damage_type_id,
  store_id,
  damage_type_name
)
VALUES (
  $1,
  $2,
  $3
);

-- name: IsDamageTypeNameTaken :one
SELECT 1
FROM damage_types
WHERE damage_types.store_id = $1 AND LOWER(damage_types.damage_type_name) = LOWER(sqlc.arg('damage_type_name'));
