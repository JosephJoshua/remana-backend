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

type SQLSalesPersonRepository struct {
	queries *gensql.Queries
}

func NewSQLSalesPersonRepository(db *pgxpool.Pool) *SQLSalesPersonRepository {
	return &SQLSalesPersonRepository{
		queries: gensql.New(db),
	}
}

func (s *SQLSalesPersonRepository) CreateSalesPerson(
	ctx context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
) error {
	if err := s.queries.CreateSalesPerson(ctx, gensql.CreateSalesPersonParams{
		SalesPersonID:   typemapper.UUIDToPgtypeUUID(id),
		StoreID:         typemapper.UUIDToPgtypeUUID(storeID),
		SalesPersonName: name,
	}); err != nil {
		return fmt.Errorf("failed to create sales person: %w", err)
	}

	return nil
}

func (s *SQLSalesPersonRepository) IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error) {
	_, err := s.queries.IsSalesPersonNameTaken(ctx, gensql.IsSalesPersonNameTakenParams{
		StoreID:         typemapper.UUIDToPgtypeUUID(storeID),
		SalesPersonName: name,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if name is taken: %w", err)
	}

	return true, nil
}
