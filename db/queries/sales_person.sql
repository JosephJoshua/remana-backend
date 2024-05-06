-- name: CreateSalesPerson :exec
INSERT INTO sales_persons (
  sales_person_id,
  store_id,
  sales_person_name
)
VALUES (
  $1,
  $2,
  $3
);

-- name: IsSalesPersonNameTaken :one
SELECT 1
FROM sales_persons
WHERE sales_persons.store_id = $1 AND LOWER(sales_persons.sales_person_name) = LOWER(sqlc.arg('sales_person_name'));
