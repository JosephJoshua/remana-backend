-- name: CreatePaymentMethod :exec
INSERT INTO payment_methods (
  payment_method_id,
  store_id,
  payment_method_name
)
VALUES (
  $1,
  $2,
  $3
);

-- name: IsPaymentMethodNameTaken :one
SELECT 1
FROM payment_methods
WHERE payment_methods.store_id = $1 AND LOWER(payment_methods.payment_method_name) = LOWER(sqlc.arg('payment_method_name'));
