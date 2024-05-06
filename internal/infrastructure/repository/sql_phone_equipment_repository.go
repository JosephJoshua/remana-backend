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

type SQLPhoneEquipmentRepository struct {
	queries *gensql.Queries
}

func NewSQLPhoneEquipmentRepository(db *pgxpool.Pool) *SQLPhoneEquipmentRepository {
	return &SQLPhoneEquipmentRepository{
		queries: gensql.New(db),
	}
}

func (s *SQLPhoneEquipmentRepository) CreatePhoneEquipment(
	ctx context.Context,
	id uuid.UUID,
	storeID uuid.UUID,
	name string,
) error {
	if err := s.queries.CreatePhoneEquipment(ctx, gensql.CreatePhoneEquipmentParams{
		PhoneEquipmentID:   typemapper.UUIDToPgtypeUUID(id),
		StoreID:            typemapper.UUIDToPgtypeUUID(storeID),
		PhoneEquipmentName: name,
	}); err != nil {
		return fmt.Errorf("failed to create phone equipment: %w", err)
	}

	return nil
}

func (s *SQLPhoneEquipmentRepository) IsNameTaken(ctx context.Context, storeID uuid.UUID, name string) (bool, error) {
	_, err := s.queries.IsPhoneEquipmentNameTaken(ctx, gensql.IsPhoneEquipmentNameTakenParams{
		StoreID:            typemapper.UUIDToPgtypeUUID(storeID),
		PhoneEquipmentName: name,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if name is taken: %w", err)
	}

	return true, nil
}
