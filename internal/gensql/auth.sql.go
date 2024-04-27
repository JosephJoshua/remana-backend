// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: auth.sql

package gensql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteLoginCodeByID = `-- name: DeleteLoginCodeByID :exec
DELETE FROM login_codes
WHERE login_codes.login_code_id = $1
`

func (q *Queries) DeleteLoginCodeByID(ctx context.Context, loginCodeID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteLoginCodeByID, loginCodeID)
	return err
}

const getLoginCodeByUserIDAndCode = `-- name: GetLoginCodeByUserIDAndCode :one
SELECT login_codes.login_code_id
FROM login_codes
WHERE login_codes.user_id = $1 AND login_codes.login_code = $2
`

type GetLoginCodeByUserIDAndCodeParams struct {
	UserID    pgtype.UUID
	LoginCode string
}

func (q *Queries) GetLoginCodeByUserIDAndCode(ctx context.Context, arg GetLoginCodeByUserIDAndCodeParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, getLoginCodeByUserIDAndCode, arg.UserID, arg.LoginCode)
	var login_code_id pgtype.UUID
	err := row.Scan(&login_code_id)
	return login_code_id, err
}

const getUserByUsernameAndStoreCode = `-- name: GetUserByUsernameAndStoreCode :one
SELECT users.user_id, users.user_password, roles.is_store_admin
FROM users
LEFT JOIN stores ON stores.store_id = users.store_id
LEFT JOIN roles ON roles.role_id = users.role_id
WHERE users.username = $1 AND stores.store_code = $2
LIMIT 1
`

type GetUserByUsernameAndStoreCodeParams struct {
	Username  string
	StoreCode string
}

type GetUserByUsernameAndStoreCodeRow struct {
	UserID       pgtype.UUID
	UserPassword string
	IsStoreAdmin pgtype.Bool
}

func (q *Queries) GetUserByUsernameAndStoreCode(ctx context.Context, arg GetUserByUsernameAndStoreCodeParams) (GetUserByUsernameAndStoreCodeRow, error) {
	row := q.db.QueryRow(ctx, getUserByUsernameAndStoreCode, arg.Username, arg.StoreCode)
	var i GetUserByUsernameAndStoreCodeRow
	err := row.Scan(&i.UserID, &i.UserPassword, &i.IsStoreAdmin)
	return i, err
}