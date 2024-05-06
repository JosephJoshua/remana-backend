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

type SQLPhoneConditionRepository struct {
	queries *gensql.Queries
}

func NewSQLPhoneConditionRepository(db *pgxpool.Pool) *SQLPhoneConditionRepository {
	return &SQLPhoneConditionRepository{
		queries: gensql.New(db),
	}
}

func (s *SQLPhoneConditionRepository) CreatePhoneCondition(
	ctx context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
) error {
	if err := s.queries.CreatePhoneCondition(ctx, gensql.CreatePhoneConditionParams{
		PhoneConditionID:   typemapper.UUIDToPgtypeUUID(id),
		StoreID:            typemapper.UUIDToPgtypeUUID(storeID),
		PhoneConditionName: name,
	}); err != nil {
		return fmt.Errorf("failed to create phone condition: %w", err)
	}

	return nil
}

func (s *SQLPhoneConditionRepository) IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error) {
	_, err := s.queries.IsPhoneConditionNameTaken(ctx, gensql.IsPhoneConditionNameTakenParams{
		StoreID:            typemapper.UUIDToPgtypeUUID(storeID),
		PhoneConditionName: name,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if name is taken: %w", err)
	}

	return true, nil
}
