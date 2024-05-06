// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: testing.sql

package gensql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getDamageTypeForTesting = `-- name: GetDamageTypeForTesting :one
SELECT
  damage_types.damage_type_id, damage_types.store_id, damage_types.damage_type_name
FROM damage_types
WHERE damage_types.damage_type_id = $1
LIMIT 1
`

func (q *Queries) GetDamageTypeForTesting(ctx context.Context, damageTypeID pgtype.UUID) (DamageType, error) {
	row := q.db.QueryRow(ctx, getDamageTypeForTesting, damageTypeID)
	var i DamageType
	err := row.Scan(&i.DamageTypeID, &i.StoreID, &i.DamageTypeName)
	return i, err
}

const getRepairOrderCostsForTesting = `-- name: GetRepairOrderCostsForTesting :many
SELECT
  repair_order_costs.repair_order_cost_id, repair_order_costs.repair_order_id, repair_order_costs.amount, repair_order_costs.reason, repair_order_costs.creation_time
FROM repair_order_costs
WHERE repair_order_costs.repair_order_id = $1
`

func (q *Queries) GetRepairOrderCostsForTesting(ctx context.Context, repairOrderID pgtype.UUID) ([]RepairOrderCost, error) {
	rows, err := q.db.Query(ctx, getRepairOrderCostsForTesting, repairOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RepairOrderCost
	for rows.Next() {
		var i RepairOrderCost
		if err := rows.Scan(
			&i.RepairOrderCostID,
			&i.RepairOrderID,
			&i.Amount,
			&i.Reason,
			&i.CreationTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRepairOrderDamagesForTesting = `-- name: GetRepairOrderDamagesForTesting :many
SELECT
  repair_order_damages.repair_order_damage_id, repair_order_damages.repair_order_id, repair_order_damages.damage_name
FROM repair_order_damages
WHERE repair_order_damages.repair_order_id = $1
`

func (q *Queries) GetRepairOrderDamagesForTesting(ctx context.Context, repairOrderID pgtype.UUID) ([]RepairOrderDamage, error) {
	rows, err := q.db.Query(ctx, getRepairOrderDamagesForTesting, repairOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RepairOrderDamage
	for rows.Next() {
		var i RepairOrderDamage
		if err := rows.Scan(&i.RepairOrderDamageID, &i.RepairOrderID, &i.DamageName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRepairOrderForTesting = `-- name: GetRepairOrderForTesting :one
SELECT
  repair_orders.repair_order_id, repair_orders.creation_time, repair_orders.slug, repair_orders.store_id, repair_orders.customer_name, repair_orders.contact_number, repair_orders.phone_type, repair_orders.imei, repair_orders.parts_not_checked_yet, repair_orders.color, repair_orders.passcode_or_pattern, repair_orders.is_pattern_locked, repair_orders.pick_up_time, repair_orders.completion_time, repair_orders.cancellation_time, repair_orders.cancellation_reason, repair_orders.confirmation_time, repair_orders.confirmation_content, repair_orders.warranty_days, repair_orders.down_payment_amount, repair_orders.down_payment_method_id, repair_orders.repayment_amount, repair_orders.repayment_method_id, repair_orders.technician_id, repair_orders.sales_person_id
FROM repair_orders
WHERE repair_orders.repair_order_id = $1
LIMIT 1
`

func (q *Queries) GetRepairOrderForTesting(ctx context.Context, repairOrderID pgtype.UUID) (RepairOrder, error) {
	row := q.db.QueryRow(ctx, getRepairOrderForTesting, repairOrderID)
	var i RepairOrder
	err := row.Scan(
		&i.RepairOrderID,
		&i.CreationTime,
		&i.Slug,
		&i.StoreID,
		&i.CustomerName,
		&i.ContactNumber,
		&i.PhoneType,
		&i.Imei,
		&i.PartsNotCheckedYet,
		&i.Color,
		&i.PasscodeOrPattern,
		&i.IsPatternLocked,
		&i.PickUpTime,
		&i.CompletionTime,
		&i.CancellationTime,
		&i.CancellationReason,
		&i.ConfirmationTime,
		&i.ConfirmationContent,
		&i.WarrantyDays,
		&i.DownPaymentAmount,
		&i.DownPaymentMethodID,
		&i.RepaymentAmount,
		&i.RepaymentMethodID,
		&i.TechnicianID,
		&i.SalesPersonID,
	)
	return i, err
}

const getRepairOrderPhoneConditionsForTesting = `-- name: GetRepairOrderPhoneConditionsForTesting :many
SELECT
  repair_order_phone_conditions.repair_order_phone_condition_id, repair_order_phone_conditions.repair_order_id, repair_order_phone_conditions.phone_condition_name
FROM repair_order_phone_conditions
WHERE repair_order_phone_conditions.repair_order_id = $1
`

func (q *Queries) GetRepairOrderPhoneConditionsForTesting(ctx context.Context, repairOrderID pgtype.UUID) ([]RepairOrderPhoneCondition, error) {
	rows, err := q.db.Query(ctx, getRepairOrderPhoneConditionsForTesting, repairOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RepairOrderPhoneCondition
	for rows.Next() {
		var i RepairOrderPhoneCondition
		if err := rows.Scan(&i.RepairOrderPhoneConditionID, &i.RepairOrderID, &i.PhoneConditionName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRepairOrderPhoneEquipmentsForTesting = `-- name: GetRepairOrderPhoneEquipmentsForTesting :many
SELECT
  repair_order_phone_equipments.repair_order_phone_equipment_id, repair_order_phone_equipments.repair_order_id, repair_order_phone_equipments.phone_equipment_name
FROM repair_order_phone_equipments
WHERE repair_order_phone_equipments.repair_order_id = $1
`

func (q *Queries) GetRepairOrderPhoneEquipmentsForTesting(ctx context.Context, repairOrderID pgtype.UUID) ([]RepairOrderPhoneEquipment, error) {
	rows, err := q.db.Query(ctx, getRepairOrderPhoneEquipmentsForTesting, repairOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RepairOrderPhoneEquipment
	for rows.Next() {
		var i RepairOrderPhoneEquipment
		if err := rows.Scan(&i.RepairOrderPhoneEquipmentID, &i.RepairOrderID, &i.PhoneEquipmentName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRepairOrderPhotosForTesting = `-- name: GetRepairOrderPhotosForTesting :many
SELECT
  repair_order_photos.repair_order_photo_id, repair_order_photos.repair_order_id, repair_order_photos.photo_url
FROM repair_order_photos
WHERE repair_order_photos.repair_order_id = $1
`

func (q *Queries) GetRepairOrderPhotosForTesting(ctx context.Context, repairOrderID pgtype.UUID) ([]RepairOrderPhoto, error) {
	rows, err := q.db.Query(ctx, getRepairOrderPhotosForTesting, repairOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RepairOrderPhoto
	for rows.Next() {
		var i RepairOrderPhoto
		if err := rows.Scan(&i.RepairOrderPhotoID, &i.RepairOrderID, &i.PhotoUrl); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSalesPersonForTesting = `-- name: GetSalesPersonForTesting :one
SELECT
  sales_persons.sales_person_id, sales_persons.store_id, sales_persons.sales_person_name
FROM sales_persons
WHERE sales_persons.sales_person_id = $1
LIMIT 1
`

func (q *Queries) GetSalesPersonForTesting(ctx context.Context, salesPersonID pgtype.UUID) (SalesPerson, error) {
	row := q.db.QueryRow(ctx, getSalesPersonForTesting, salesPersonID)
	var i SalesPerson
	err := row.Scan(&i.SalesPersonID, &i.StoreID, &i.SalesPersonName)
	return i, err
}

const getTechnicianForTesting = `-- name: GetTechnicianForTesting :one
SELECT
  technicians.technician_id, technicians.store_id, technicians.technician_name
FROM technicians
WHERE technicians.technician_id = $1
LIMIT 1
`

func (q *Queries) GetTechnicianForTesting(ctx context.Context, technicianID pgtype.UUID) (Technician, error) {
	row := q.db.QueryRow(ctx, getTechnicianForTesting, technicianID)
	var i Technician
	err := row.Scan(&i.TechnicianID, &i.StoreID, &i.TechnicianName)
	return i, err
}

const seedDamageType = `-- name: SeedDamageType :one
INSERT INTO damage_types (damage_type_id, damage_type_name, store_id)
VALUES ($1, $2, $3)
RETURNING damage_type_id
`

type SeedDamageTypeParams struct {
	DamageTypeID   pgtype.UUID
	DamageTypeName string
	StoreID        pgtype.UUID
}

func (q *Queries) SeedDamageType(ctx context.Context, arg SeedDamageTypeParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedDamageType, arg.DamageTypeID, arg.DamageTypeName, arg.StoreID)
	var damage_type_id pgtype.UUID
	err := row.Scan(&damage_type_id)
	return damage_type_id, err
}

const seedLoginCode = `-- name: SeedLoginCode :one
INSERT INTO login_codes (login_code_id, user_id, login_code)
VALUES ($1, $2, $3)
RETURNING login_code_id
`

type SeedLoginCodeParams struct {
	LoginCodeID pgtype.UUID
	UserID      pgtype.UUID
	LoginCode   string
}

func (q *Queries) SeedLoginCode(ctx context.Context, arg SeedLoginCodeParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedLoginCode, arg.LoginCodeID, arg.UserID, arg.LoginCode)
	var login_code_id pgtype.UUID
	err := row.Scan(&login_code_id)
	return login_code_id, err
}

const seedPaymentMethod = `-- name: SeedPaymentMethod :one
INSERT INTO payment_methods (payment_method_id, payment_method_name, store_id)
VALUES ($1, $2, $3)
RETURNING payment_method_id
`

type SeedPaymentMethodParams struct {
	PaymentMethodID   pgtype.UUID
	PaymentMethodName string
	StoreID           pgtype.UUID
}

func (q *Queries) SeedPaymentMethod(ctx context.Context, arg SeedPaymentMethodParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedPaymentMethod, arg.PaymentMethodID, arg.PaymentMethodName, arg.StoreID)
	var payment_method_id pgtype.UUID
	err := row.Scan(&payment_method_id)
	return payment_method_id, err
}

const seedPhoneCondition = `-- name: SeedPhoneCondition :one
INSERT INTO phone_conditions (phone_condition_id, phone_condition_name, store_id)
VALUES ($1, $2, $3)
RETURNING phone_condition_id
`

type SeedPhoneConditionParams struct {
	PhoneConditionID   pgtype.UUID
	PhoneConditionName string
	StoreID            pgtype.UUID
}

func (q *Queries) SeedPhoneCondition(ctx context.Context, arg SeedPhoneConditionParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedPhoneCondition, arg.PhoneConditionID, arg.PhoneConditionName, arg.StoreID)
	var phone_condition_id pgtype.UUID
	err := row.Scan(&phone_condition_id)
	return phone_condition_id, err
}

const seedPhoneEquipment = `-- name: SeedPhoneEquipment :one
INSERT INTO phone_equipments (phone_equipment_id, phone_equipment_name, store_id)
VALUES ($1, $2, $3)
RETURNING phone_equipment_id
`

type SeedPhoneEquipmentParams struct {
	PhoneEquipmentID   pgtype.UUID
	PhoneEquipmentName string
	StoreID            pgtype.UUID
}

func (q *Queries) SeedPhoneEquipment(ctx context.Context, arg SeedPhoneEquipmentParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedPhoneEquipment, arg.PhoneEquipmentID, arg.PhoneEquipmentName, arg.StoreID)
	var phone_equipment_id pgtype.UUID
	err := row.Scan(&phone_equipment_id)
	return phone_equipment_id, err
}

const seedRole = `-- name: SeedRole :one
INSERT INTO roles (role_id, role_name, store_id, is_store_admin)
VALUES ($1, $2, $3, $4)
RETURNING role_id
`

type SeedRoleParams struct {
	RoleID       pgtype.UUID
	RoleName     string
	StoreID      pgtype.UUID
	IsStoreAdmin bool
}

func (q *Queries) SeedRole(ctx context.Context, arg SeedRoleParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedRole,
		arg.RoleID,
		arg.RoleName,
		arg.StoreID,
		arg.IsStoreAdmin,
	)
	var role_id pgtype.UUID
	err := row.Scan(&role_id)
	return role_id, err
}

const seedSalesPerson = `-- name: SeedSalesPerson :one
INSERT INTO sales_persons (sales_person_id, sales_person_name, store_id)
VALUES ($1, $2, $3)
RETURNING sales_person_id
`

type SeedSalesPersonParams struct {
	SalesPersonID   pgtype.UUID
	SalesPersonName string
	StoreID         pgtype.UUID
}

func (q *Queries) SeedSalesPerson(ctx context.Context, arg SeedSalesPersonParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedSalesPerson, arg.SalesPersonID, arg.SalesPersonName, arg.StoreID)
	var sales_person_id pgtype.UUID
	err := row.Scan(&sales_person_id)
	return sales_person_id, err
}

const seedStore = `-- name: SeedStore :one
INSERT INTO stores (store_id, store_name, store_code, store_address, phone_number)
VALUES ($1, $2, $3, $4, $5)
RETURNING store_id
`

type SeedStoreParams struct {
	StoreID      pgtype.UUID
	StoreName    string
	StoreCode    string
	StoreAddress string
	PhoneNumber  string
}

func (q *Queries) SeedStore(ctx context.Context, arg SeedStoreParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedStore,
		arg.StoreID,
		arg.StoreName,
		arg.StoreCode,
		arg.StoreAddress,
		arg.PhoneNumber,
	)
	var store_id pgtype.UUID
	err := row.Scan(&store_id)
	return store_id, err
}

const seedTechnician = `-- name: SeedTechnician :one
INSERT INTO technicians (technician_id, technician_name, store_id)
VALUES ($1, $2, $3)
RETURNING technician_id
`

type SeedTechnicianParams struct {
	TechnicianID   pgtype.UUID
	TechnicianName string
	StoreID        pgtype.UUID
}

func (q *Queries) SeedTechnician(ctx context.Context, arg SeedTechnicianParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedTechnician, arg.TechnicianID, arg.TechnicianName, arg.StoreID)
	var technician_id pgtype.UUID
	err := row.Scan(&technician_id)
	return technician_id, err
}

const seedUser = `-- name: SeedUser :one
INSERT INTO users (user_id, username, user_password, role_id, store_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id
`

type SeedUserParams struct {
	UserID       pgtype.UUID
	Username     string
	UserPassword string
	RoleID       pgtype.UUID
	StoreID      pgtype.UUID
}

func (q *Queries) SeedUser(ctx context.Context, arg SeedUserParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, seedUser,
		arg.UserID,
		arg.Username,
		arg.UserPassword,
		arg.RoleID,
		arg.StoreID,
	)
	var user_id pgtype.UUID
	err := row.Scan(&user_id)
	return user_id, err
}
