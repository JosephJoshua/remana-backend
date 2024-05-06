-- name: CreatePhoneCondition :exec
INSERT INTO phone_conditions (
  phone_condition_id,
  store_id,
  phone_condition_name
)
VALUES (
  $1,
  $2,
  $3
);

-- name: IsPhoneConditionNameTaken :one
SELECT 1
FROM phone_conditions
WHERE phone_conditions.store_id = $1 AND LOWER(phone_conditions.phone_condition_name) = LOWER(sqlc.arg('phone_condition_name'));
