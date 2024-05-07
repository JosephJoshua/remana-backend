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

type SQLPaymentMethodRepository struct {
	queries *gensql.Queries
}

func NewSQLPaymentMethodRepository(db *pgxpool.Pool) *SQLPaymentMethodRepository {
	return &SQLPaymentMethodRepository{
		queries: gensql.New(db),
	}
}

func (s *SQLPaymentMethodRepository) CreatePaymentMethod(
	ctx context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
) error {
	if err := s.queries.CreatePaymentMethod(ctx, gensql.CreatePaymentMethodParams{
		PaymentMethodID:   typemapper.UUIDToPgtypeUUID(id),
		StoreID:           typemapper.UUIDToPgtypeUUID(storeID),
		PaymentMethodName: name,
	}); err != nil {
		return fmt.Errorf("failed to create payment method: %w", err)
	}

	return nil
}

func (s *SQLPaymentMethodRepository) IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error) {
	_, err := s.queries.IsPaymentMethodNameTaken(ctx, gensql.IsPaymentMethodNameTakenParams{
		StoreID:           typemapper.UUIDToPgtypeUUID(storeID),
		PaymentMethodName: name,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if name is taken: %w", err)
	}

	return true, nil
}
