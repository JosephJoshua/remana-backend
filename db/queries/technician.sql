-- name: CreateTechnician :exec
INSERT INTO technicians (
  technician_id,
  store_id,
  technician_name
)
VALUES (
  $1,
  $2,
  $3
);

-- name: IsTechnicianNameTaken :one
SELECT 1
FROM technicians
WHERE technicians.store_id = $1 AND LOWER(technicians.technician_name) = LOWER(sqlc.arg('technician_name'));
