-- name: CreatePhoneEquipment :exec
INSERT INTO phone_equipments (
  phone_equipment_id,
  store_id,
  phone_equipment_name
)
VALUES (
  $1,
  $2,
  $3
);

-- name: IsPhoneEquipmentNameTaken :one
SELECT 1
FROM phone_equipments
WHERE phone_equipments.store_id = $1 AND phone_equipments.phone_equipment_name = $2;
