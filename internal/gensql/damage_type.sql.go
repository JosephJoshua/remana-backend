// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: damage_type.sql

package gensql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createDamageType = `-- name: CreateDamageType :exec
INSERT INTO damage_types (
  damage_type_id,
  store_id,
  damage_type_name
)
VALUES (
  $1,
  $2,
  $3
)
`

type CreateDamageTypeParams struct {
	DamageTypeID   pgtype.UUID
	StoreID        pgtype.UUID
	DamageTypeName string
}

func (q *Queries) CreateDamageType(ctx context.Context, arg CreateDamageTypeParams) error {
	_, err := q.db.Exec(ctx, createDamageType, arg.DamageTypeID, arg.StoreID, arg.DamageTypeName)
	return err
}

const isDamageTypeNameTaken = `-- name: IsDamageTypeNameTaken :one
SELECT 1
FROM damage_types
WHERE damage_types.store_id = $1 AND damage_types.damage_type_name = $2
`

type IsDamageTypeNameTakenParams struct {
	StoreID        pgtype.UUID
	DamageTypeName string
}

func (q *Queries) IsDamageTypeNameTaken(ctx context.Context, arg IsDamageTypeNameTakenParams) (int32, error) {
	row := q.db.QueryRow(ctx, isDamageTypeNameTaken, arg.StoreID, arg.DamageTypeName)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}