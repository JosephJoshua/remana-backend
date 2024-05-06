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

type SQLDamageTypeRepository struct {
	queries *gensql.Queries
}

func NewSQLDamageTypeRepository(db *pgxpool.Pool) *SQLDamageTypeRepository {
	return &SQLDamageTypeRepository{
		queries: gensql.New(db),
	}
}

func (s *SQLDamageTypeRepository) CreateDamageType(
	ctx context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
) error {
	if err := s.queries.CreateDamageType(ctx, gensql.CreateDamageTypeParams{
		DamageTypeID:   typemapper.UUIDToPgtypeUUID(id),
		StoreID:        typemapper.UUIDToPgtypeUUID(storeID),
		DamageTypeName: name,
	}); err != nil {
		return fmt.Errorf("failed to create damage type: %w", err)
	}

	return nil
}

func (s *SQLDamageTypeRepository) IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error) {
	_, err := s.queries.IsDamageTypeNameTaken(ctx, gensql.IsDamageTypeNameTakenParams{
		StoreID:        typemapper.UUIDToPgtypeUUID(storeID),
		DamageTypeName: name,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if name is taken: %w", err)
	}

	return true, nil
}
