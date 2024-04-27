-- name: CreateRepairOrder :exec
INSERT INTO repair_orders (
  repair_order_id,
  creation_time,
  slug,
  store_id,
  customer_name,
  contact_number,
  phone_type,
  color,
  sales_id,
  technician_id,
  imei,
  parts_not_checked_yet,
  passcode_or_pattern,
  is_pattern_locked,
  down_payment_amount,
  down_payment_method_id
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  $15,
  $16
);

-- name: AddDamagesToRepairOrder :copyfrom
INSERT INTO repair_order_damages (
  repair_order_damage_id,
  repair_order_id,
  damage_name
) VALUES (
  $1,
  $2,
  $3
);

-- name: AddPhoneConditionsToRepairOrder :copyfrom
INSERT INTO repair_order_phone_conditions (
  repair_order_phone_condition_id,
  repair_order_id,
  phone_condition_name
) VALUES (
  $1,
  $2,
  $3
);

-- name: AddPhoneEquipmentsToRepairOrder :copyfrom
INSERT INTO repair_order_phone_equipments (
  repair_order_phone_equipment_id,
  repair_order_id,
  phone_equipment_name
) VALUES (
  $1,
  $2,
  $3
);

-- name: AddPhotosToRepairOrder :copyfrom
INSERT INTO repair_order_photos (
  repair_order_photo_id,
  repair_order_id,
  photo_url
) VALUES (
  $1,
  $2,
  $3
);

-- name: AddCostsToRepairOrder :copyfrom
INSERT INTO repair_order_costs (
  repair_order_cost_id,
  repair_order_id,
  amount,
  reason,
  creation_time
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
);

-- name: DoesSalesExist :one
SELECT 1
FROM sales
WHERE sales.store_id = $1 AND sales.sales_id = $2;

-- name: DoesTechnicianExist :one
SELECT 1
FROM technicians
WHERE technicians.store_id = $1 AND technicians.technician_id = $2;

-- name: DoesPaymentMethodExist :one
SELECT 1
FROM payment_methods
WHERE payment_methods.store_id = $1 AND payment_methods.payment_method_id = $2;

-- name: GetDamageNamesByIDs :many
SELECT damage_types.damage_type_name
FROM damage_types
WHERE damage_types.store_id = $1 AND damage_types.damage_type_id = ANY(sqlc.arg(ids)::UUID[]);

-- name: GetPhoneConditionNamesByIDs :many
SELECT phone_conditions.phone_condition_name
FROM phone_conditions
WHERE phone_conditions.store_id = $1 AND phone_conditions.phone_condition_id = ANY(sqlc.arg(ids)::UUID[]);

-- name: GetPhoneEquipmentNamesByIDs :many
SELECT phone_equipments.phone_equipment_name
FROM phone_equipments
WHERE phone_equipments.store_id = $1 AND phone_equipments.phone_equipment_id = ANY(sqlc.arg(ids)::UUID[]);

-- name: IsRepairOrderSlugTaken :one
SELECT 1
FROM repair_orders
WHERE repair_orders.store_id = $1 AND repair_orders.slug = $2;
