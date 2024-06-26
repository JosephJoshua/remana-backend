// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: sales_person.sql

package gensql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSalesPerson = `-- name: CreateSalesPerson :exec
INSERT INTO sales_persons (
  sales_person_id,
  store_id,
  sales_person_name
)
VALUES (
  $1,
  $2,
  $3
)
`

type CreateSalesPersonParams struct {
	SalesPersonID   pgtype.UUID
	StoreID         pgtype.UUID
	SalesPersonName string
}

func (q *Queries) CreateSalesPerson(ctx context.Context, arg CreateSalesPersonParams) error {
	_, err := q.db.Exec(ctx, createSalesPerson, arg.SalesPersonID, arg.StoreID, arg.SalesPersonName)
	return err
}

const isSalesPersonNameTaken = `-- name: IsSalesPersonNameTaken :one
SELECT 1
FROM sales_persons
WHERE sales_persons.store_id = $1 AND LOWER(sales_persons.sales_person_name) = LOWER($2)
`

type IsSalesPersonNameTakenParams struct {
	StoreID         pgtype.UUID
	SalesPersonName string
}

func (q *Queries) IsSalesPersonNameTaken(ctx context.Context, arg IsSalesPersonNameTakenParams) (int32, error) {
	row := q.db.QueryRow(ctx, isSalesPersonNameTaken, arg.StoreID, arg.SalesPersonName)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}
