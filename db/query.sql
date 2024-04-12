-- name: GetStoreByID :one
SELECT * FROM stores
WHERE store_id = $1 LIMIT 1;
