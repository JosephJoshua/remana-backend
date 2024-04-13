package repository

import (
	"context"
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLUserRepository struct {
	queries *gensql.Queries
}

func NewSQLUserRepository(db *pgxpool.Pool) *SQLUserRepository {
	return &SQLUserRepository{
		queries: gensql.New(db),
	}
}

func (r *SQLUserRepository) GetUserDetailsByID(ctx context.Context, userID uuid.UUID) (readmodel.UserDetails, error) {
	user, err := r.queries.GetUserDetailsByID(ctx, pgtype.UUID{
		Bytes: userID,
		Valid: true,
	})

	var emptyUser readmodel.UserDetails
	if err != nil {
		return emptyUser, fmt.Errorf("failed to get user details by ID: %w", err)
	}

	userID, err = pgtypeUUIDToGoogleUUID(user.UserID)
	if err != nil {
		return emptyUser, fmt.Errorf("failed to parse user ID from bytes: %w", err)
	}

	roleID, err := pgtypeUUIDToGoogleUUID(user.RoleID)
	if err != nil {
		return emptyUser, fmt.Errorf("failed to parse role ID from bytes: %w", err)
	}

	storeID, err := pgtypeUUIDToGoogleUUID(user.StoreID)
	if err != nil {
		return emptyUser, fmt.Errorf("failed to parse store ID from bytes: %w", err)
	}

	return readmodel.UserDetails{
		ID:       userID,
		Username: user.Username,
		Role: readmodel.UserDetailsRole{
			ID:           roleID,
			Name:         user.RoleName.String,
			IsStoreAdmin: user.IsStoreAdmin.Bool,
		},
		Store: readmodel.UserDetailsStore{
			ID:   storeID,
			Name: user.StoreName.String,
			Code: user.StoreCode.String,
		},
	}, nil
}