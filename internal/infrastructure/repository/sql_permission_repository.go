package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLPermissionRepository struct {
	queries *gensql.Queries
	db      *pgxpool.Pool
}

func NewSQLPermissionRepository(db *pgxpool.Pool) *SQLPermissionRepository {
	return &SQLPermissionRepository{
		queries: gensql.New(db),
		db:      db,
	}
}

func (s *SQLPermissionRepository) CreateRole(
	ctx context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
	isStoreAdmin bool,
) error {
	if err := s.queries.CreateRole(ctx, gensql.CreateRoleParams{
		RoleID:       typemapper.UUIDToPgtypeUUID(id),
		StoreID:      typemapper.UUIDToPgtypeUUID(storeID),
		RoleName:     name,
		IsStoreAdmin: isStoreAdmin,
	}); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (s *SQLPermissionRepository) IsRoleNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error) {
	_, err := s.queries.IsRoleNameTaken(ctx, gensql.IsRoleNameTakenParams{
		StoreID:  typemapper.UUIDToPgtypeUUID(storeID),
		RoleName: name,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if name is taken: %w", err)
	}

	return true, nil
}

func (s *SQLPermissionRepository) GetPermissionIDs(
	ctx context.Context,
	permissions []permission.GetPermissionIDDetail,
) ([]uuid.UUID, error) {
	query := `
		SELECT permissions.permission_id
		FROM JSON_TO_RECORDSET($1::JSON) AS input (group_name TEXT, permission_name TEXT)
		LEFT JOIN permission_groups ON
			permission_groups.permission_group_name = input.group_name
		LEFT JOIN permissions ON
			permissions.permission_name = input.permission_name AND
			permissions.permission_group_id = permission_groups.permission_group_id
		WHERE
			permissions.permission_id IS NOT NULL
	`

	permissionParams := make([]map[string]interface{}, 0, len(permissions))
	for _, p := range permissions {
		permissionParams = append(permissionParams, map[string]interface{}{
			"group_name":      p.GroupName,
			"permission_name": p.Name,
		})
	}

	rows, err := s.db.Query(ctx, query, permissionParams)
	if errors.Is(err, pgx.ErrNoRows) {
		return []uuid.UUID{}, apperror.ErrPermissionNotFound
	} else if err != nil {
		return []uuid.UUID{}, fmt.Errorf("failed to check execute query: %w", err)
	}
	defer rows.Close()

	var items []pgtype.UUID
	for rows.Next() {
		var i pgtype.UUID
		if err = rows.Scan(&i); err != nil {
			return []uuid.UUID{}, fmt.Errorf("failed to scan row: %w", err)
		}

		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return []uuid.UUID{}, fmt.Errorf("failed to read row: %w", err)
	}

	if len(items) < len(permissions) {
		return []uuid.UUID{}, apperror.ErrPermissionNotFound
	}

	return typemapper.MustPgtypeUUIDsToUUIDs(items), nil
}

func (s *SQLPermissionRepository) DoesRoleExist(ctx context.Context, roleID uuid.UUID) (bool, error) {
	_, err := s.queries.DoesRoleExist(ctx, typemapper.UUIDToPgtypeUUID(roleID))

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if role exists: %w", err)
	}

	return true, nil
}

func (s *SQLPermissionRepository) AssignPermissionsToRole(
	ctx context.Context,
	roleID uuid.UUID,
	permissionIDs []uuid.UUID,
) error {
	params := make([]gensql.AssignPermissionsToRoleParams, 0, len(permissionIDs))
	for _, id := range permissionIDs {
		params = append(params, gensql.AssignPermissionsToRoleParams{
			RoleID:       typemapper.UUIDToPgtypeUUID(roleID),
			PermissionID: typemapper.UUIDToPgtypeUUID(id),
		})
	}

	if _, err := s.queries.AssignPermissionsToRole(ctx, params); err != nil {
		return fmt.Errorf("failed to assign permissions to role: %w", err)
	}

	return nil
}
