package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLPermissionRepository struct {
	queries *gensql.Queries
}

func NewSQLPermissionRepository(db *pgxpool.Pool) *SQLPermissionRepository {
	return &SQLPermissionRepository{
		queries: gensql.New(db),
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
