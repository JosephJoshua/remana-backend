package repository

import (
	"context"
	"slices"

	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
)

type MemoryAuthRepository struct {
	users      []domain.User
	loginCodes []domain.LoginCode
}

func NewMemoryAuthRepository(users []domain.User, loginCodes []domain.LoginCode) *MemoryAuthRepository {
	return &MemoryAuthRepository{users: users, loginCodes: loginCodes}
}

func (r *MemoryAuthRepository) GetUserByUsernameAndStoreCode(
	_ context.Context,
	username string,
	storeCode string,
) (domain.User, error) {
	i := slices.IndexFunc(r.users, func(u domain.User) bool {
		return u.Username() == username && u.StoreCode() == storeCode
	})

	if i == -1 {
		return nil, apperror.ErrUserNotFound
	}

	return r.users[i], nil
}

func (r *MemoryAuthRepository) CheckUserLoginCode(
	_ context.Context,
	userID uuid.UUID,
	loginCode string,
) error {
	i := slices.IndexFunc(r.loginCodes, func(c domain.LoginCode) bool {
		return c.UserID() == userID && c.Code() == loginCode
	})

	if i == -1 {
		return apperror.ErrLoginCodeMismatch
	}

	return nil
}
