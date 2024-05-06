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
	var emptyUser readmodel.User

	user, err := r.queries.GetUserByUsernameAndStoreCode(ctx, gensql.GetUserByUsernameAndStoreCodeParams{
		Username:  username,
		StoreCode: storeCode,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return emptyUser, apperror.ErrUserNotFound
	} else if err != nil {
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
	if errors.Is(err, pgx.ErrNoRows) {
		return apperror.ErrLoginCodeMismatch
	} else if err != nil {
		return fmt.Errorf("failed to get login code by user ID and code: %w", err)
	}

	err = r.queries.DeleteLoginCodeByID(ctx, loginCodeID)
	if err != nil {
		return fmt.Errorf("failed to delete login code by ID: %w", err)
	}

	return nil
}

func (r *SQLAuthRepository) GetUserDetailsByID(ctx context.Context, userID uuid.UUID) (readmodel.UserDetails, error) {
	var emptyUser readmodel.UserDetails

	user, err := r.queries.GetUserDetailsByID(ctx, typemapper.UUIDToPgtypeUUID(userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return emptyUser, apperror.ErrUserNotFound
	} else if err != nil {
		return emptyUser, fmt.Errorf("failed to get user details by ID: %w", err)
	}

	userID, err = typemapper.PgtypeUUIDToUUID(user.UserID)
	if err != nil {
		return emptyUser, fmt.Errorf("failed to parse user ID from bytes: %w", err)
	}

	roleID, err := typemapper.PgtypeUUIDToUUID(user.RoleID)
	if err != nil {
		return emptyUser, fmt.Errorf("failed to parse role ID from bytes: %w", err)
	}

	storeID, err := typemapper.PgtypeUUIDToUUID(user.StoreID)
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
