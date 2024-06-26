// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: copyfrom.go

package gensql

import (
	"context"
)

// iteratorForAddCostsToRepairOrder implements pgx.CopyFromSource.
type iteratorForAddCostsToRepairOrder struct {
	rows                 []AddCostsToRepairOrderParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddCostsToRepairOrder) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddCostsToRepairOrder) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RepairOrderCostID,
		r.rows[0].RepairOrderID,
		r.rows[0].Amount,
		r.rows[0].Reason,
		r.rows[0].CreationTime,
	}, nil
}

func (r iteratorForAddCostsToRepairOrder) Err() error {
	return nil
}

func (q *Queries) AddCostsToRepairOrder(ctx context.Context, arg []AddCostsToRepairOrderParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"repair_order_costs"}, []string{"repair_order_cost_id", "repair_order_id", "amount", "reason", "creation_time"}, &iteratorForAddCostsToRepairOrder{rows: arg})
}

// iteratorForAddDamagesToRepairOrder implements pgx.CopyFromSource.
type iteratorForAddDamagesToRepairOrder struct {
	rows                 []AddDamagesToRepairOrderParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddDamagesToRepairOrder) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddDamagesToRepairOrder) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RepairOrderDamageID,
		r.rows[0].RepairOrderID,
		r.rows[0].DamageName,
	}, nil
}

func (r iteratorForAddDamagesToRepairOrder) Err() error {
	return nil
}

func (q *Queries) AddDamagesToRepairOrder(ctx context.Context, arg []AddDamagesToRepairOrderParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"repair_order_damages"}, []string{"repair_order_damage_id", "repair_order_id", "damage_name"}, &iteratorForAddDamagesToRepairOrder{rows: arg})
}

// iteratorForAddPhoneConditionsToRepairOrder implements pgx.CopyFromSource.
type iteratorForAddPhoneConditionsToRepairOrder struct {
	rows                 []AddPhoneConditionsToRepairOrderParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddPhoneConditionsToRepairOrder) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddPhoneConditionsToRepairOrder) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RepairOrderPhoneConditionID,
		r.rows[0].RepairOrderID,
		r.rows[0].PhoneConditionName,
	}, nil
}

func (r iteratorForAddPhoneConditionsToRepairOrder) Err() error {
	return nil
}

func (q *Queries) AddPhoneConditionsToRepairOrder(ctx context.Context, arg []AddPhoneConditionsToRepairOrderParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"repair_order_phone_conditions"}, []string{"repair_order_phone_condition_id", "repair_order_id", "phone_condition_name"}, &iteratorForAddPhoneConditionsToRepairOrder{rows: arg})
}

// iteratorForAddPhoneEquipmentsToRepairOrder implements pgx.CopyFromSource.
type iteratorForAddPhoneEquipmentsToRepairOrder struct {
	rows                 []AddPhoneEquipmentsToRepairOrderParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddPhoneEquipmentsToRepairOrder) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddPhoneEquipmentsToRepairOrder) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RepairOrderPhoneEquipmentID,
		r.rows[0].RepairOrderID,
		r.rows[0].PhoneEquipmentName,
	}, nil
}

func (r iteratorForAddPhoneEquipmentsToRepairOrder) Err() error {
	return nil
}

func (q *Queries) AddPhoneEquipmentsToRepairOrder(ctx context.Context, arg []AddPhoneEquipmentsToRepairOrderParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"repair_order_phone_equipments"}, []string{"repair_order_phone_equipment_id", "repair_order_id", "phone_equipment_name"}, &iteratorForAddPhoneEquipmentsToRepairOrder{rows: arg})
}

// iteratorForAddPhotosToRepairOrder implements pgx.CopyFromSource.
type iteratorForAddPhotosToRepairOrder struct {
	rows                 []AddPhotosToRepairOrderParams
	skippedFirstNextCall bool
}

func (r *iteratorForAddPhotosToRepairOrder) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAddPhotosToRepairOrder) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RepairOrderPhotoID,
		r.rows[0].RepairOrderID,
		r.rows[0].PhotoUrl,
	}, nil
}

func (r iteratorForAddPhotosToRepairOrder) Err() error {
	return nil
}

func (q *Queries) AddPhotosToRepairOrder(ctx context.Context, arg []AddPhotosToRepairOrderParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"repair_order_photos"}, []string{"repair_order_photo_id", "repair_order_id", "photo_url"}, &iteratorForAddPhotosToRepairOrder{rows: arg})
}

// iteratorForAssignPermissionsToRole implements pgx.CopyFromSource.
type iteratorForAssignPermissionsToRole struct {
	rows                 []AssignPermissionsToRoleParams
	skippedFirstNextCall bool
}

func (r *iteratorForAssignPermissionsToRole) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForAssignPermissionsToRole) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].RoleID,
		r.rows[0].PermissionID,
	}, nil
}

func (r iteratorForAssignPermissionsToRole) Err() error {
	return nil
}

func (q *Queries) AssignPermissionsToRole(ctx context.Context, arg []AssignPermissionsToRoleParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"role_permissions"}, []string{"role_id", "permission_id"}, &iteratorForAssignPermissionsToRole{rows: arg})
}
