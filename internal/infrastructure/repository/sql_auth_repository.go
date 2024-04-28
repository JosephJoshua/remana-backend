package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/JosephJoshua/remana-backend/internal/apperror"
	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLAuthRepository struct {
	queries *gensql.Queries
}

func NewSQLAuthRepository(db *pgxpool.Pool) *SQLAuthRepository {
	return &SQLAuthRepository{
		queries: gensql.New(db),
	}
}

func (r *SQLAuthRepository) GetUserByUsernameAndStoreCode(
	ctx context.Context,
	username string,
	storeCode string,
) (readmodel.User, error) {
	user, err := r.queries.GetUserByUsernameAndStoreCode(ctx, gensql.GetUserByUsernameAndStoreCodeParams{
		Username:  username,
		StoreCode: storeCode,
	})

	var emptyUser readmodel.User
	if err != nil {
		return emptyUser, fmt.Errorf("failed to get user by username and store code: %w", err)
	}

	userID, err := typemapper.PgtypeUUIDToUUID(user.UserID)
	if err != nil {
		return emptyUser, fmt.Errorf("failed to parse user ID from bytes: %w", err)
	}

	return readmodel.User{
		ID:           userID,
		Password:     user.UserPassword,
		IsStoreAdmin: user.IsStoreAdmin.Bool,
	}, nil
}

func (r *SQLAuthRepository) CheckAndDeleteUserLoginCode(
	ctx context.Context,
	userID uuid.UUID,
	loginCode string,
) error {
	loginCodeID, err := r.queries.GetLoginCodeByUserIDAndCode(ctx, gensql.GetLoginCodeByUserIDAndCodeParams{
		UserID:    typemapper.UUIDToPgtypeUUID(userID),
		LoginCode: loginCode,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrLoginCodeMismatch
		}

		return fmt.Errorf("failed to get login code by user ID and code: %w", err)
	}

	err = r.queries.DeleteLoginCodeByID(ctx, loginCodeID)
	if err != nil {
		return fmt.Errorf("failed to delete login code by ID: %w", err)
	}

	return nil
}
