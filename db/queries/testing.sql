-- name: SeedStore :one
INSERT INTO stores (store_id, store_name, store_code, store_address, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING store_id;

-- name: SeedRole :one
INSERT INTO roles (role_id, role_name, store_id, is_store_admin)
VALUES ($1, $2, $3, $4)
RETURNING role_id;

-- name: SeedUser :one
INSERT INTO users (user_id, username, user_password, role_id, store_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id;

-- name: SeedLoginCode :one
INSERT INTO login_codes (login_code_id, user_id, login_code)
VALUES ($1, $2, $3)
RETURNING login_code_id;

-- name: SeedTechnician :one
INSERT INTO technicians (technician_id, technician_name, store_id)
VALUES ($1, $2, $3)
RETURNING technician_id;

-- name: SeedSalesPerson :one
INSERT INTO sales_persons (sales_person_id, sales_person_name, store_id)
VALUES ($1, $2, $3)
RETURNING sales_person_id;

-- name: SeedDamageType :one
INSERT INTO damage_types (damage_type_id, damage_type_name, store_id)
VALUES ($1, $2, $3)
RETURNING damage_type_id;

-- name: SeedPhoneCondition :one
INSERT INTO phone_conditions (phone_condition_id, phone_condition_name, store_id)
VALUES ($1, $2, $3)
RETURNING phone_condition_id;

-- name: SeedPhoneEquipment :one
INSERT INTO phone_equipments (phone_equipment_id, phone_equipment_name, store_id)
VALUES ($1, $2, $3)
RETURNING phone_equipment_id;

-- name: SeedPaymentMethod :one
INSERT INTO payment_methods (payment_method_id, payment_method_name, store_id)
VALUES ($1, $2, $3)
RETURNING payment_method_id;

-- name: GetRoleForTesting :one
SELECT
  roles.*
FROM roles
WHERE roles.role_id = $1
LIMIT 1;

-- name: GetTechnicianForTesting :one
SELECT
  technicians.*
FROM technicians
WHERE technicians.technician_id = $1
LIMIT 1;

-- name: GetSalesPersonForTesting :one
SELECT
  sales_persons.*
FROM sales_persons
WHERE sales_persons.sales_person_id = $1
LIMIT 1;

-- name: GetDamageTypeForTesting :one
SELECT
  damage_types.*
FROM damage_types
WHERE damage_types.damage_type_id = $1
LIMIT 1;

-- name: GetPhoneConditionForTesting :one
SELECT
  phone_conditions.*
FROM phone_conditions
WHERE phone_conditions.phone_condition_id = $1
LIMIT 1;

-- name: GetPhoneEquipmentForTesting :one
SELECT
  phone_equipments.*
FROM phone_equipments
WHERE phone_equipments.phone_equipment_id = $1
LIMIT 1;

-- name: GetPaymentMethodForTesting :one
SELECT
  payment_methods.*
FROM payment_methods
WHERE payment_methods.payment_method_id = $1
LIMIT 1;

-- name: GetRepairOrderForTesting :one
SELECT
  repair_orders.*
FROM repair_orders
WHERE repair_orders.repair_order_id = $1
LIMIT 1;

-- name: GetRepairOrderDamagesForTesting :many
SELECT
  repair_order_damages.*
FROM repair_order_damages
WHERE repair_order_damages.repair_order_id = $1;

-- name: GetRepairOrderPhoneConditionsForTesting :many
SELECT
  repair_order_phone_conditions.*
FROM repair_order_phone_conditions
WHERE repair_order_phone_conditions.repair_order_id = $1;

-- name: GetRepairOrderPhoneEquipmentsForTesting :many
SELECT
  repair_order_phone_equipments.*
FROM repair_order_phone_equipments
WHERE repair_order_phone_equipments.repair_order_id = $1;

-- name: GetRepairOrderCostsForTesting :many
SELECT
  repair_order_costs.*
FROM repair_order_costs
WHERE repair_order_costs.repair_order_id = $1;

-- name: GetRepairOrderPhotosForTesting :many
SELECT
  repair_order_photos.*
FROM repair_order_photos
WHERE repair_order_photos.repair_order_id = $1;
