package repository

import (
	"context"
	"slices"

	"github.com/JosephJoshua/remana-backend/internal/shared/apperror"
	"github.com/JosephJoshua/remana-backend/internal/shared/domain"
	"github.com/google/uuid"
)

type MemoryAuthRepository struct {
	Users      []domain.User
	LoginCodes []domain.LoginCode
}

func NewMemoryAuthRepository(users []domain.User, loginCodes []domain.LoginCode) *MemoryAuthRepository {
	return &MemoryAuthRepository{Users: users, LoginCodes: loginCodes}
}

func (r *MemoryAuthRepository) GetUserByUsernameAndStoreCode(
	_ context.Context,
	username string,
	storeCode string,
) (domain.User, error) {
	i := slices.IndexFunc(r.Users, func(u domain.User) bool {
		return u.Username() == username && u.StoreCode() == storeCode
	})

	if i == -1 {
		return nil, apperror.ErrUserNotFound
	}

	return r.Users[i], nil
}

func (r *MemoryAuthRepository) CheckAndDeleteUserLoginCode(
	_ context.Context,
	userID uuid.UUID,
	loginCode string,
) error {
	i := slices.IndexFunc(r.LoginCodes, func(c domain.LoginCode) bool {
		return c.UserID() == userID && c.Code() == loginCode
	})

	if i == -1 {
		return apperror.ErrLoginCodeMismatch
	}

	r.LoginCodes[i] = r.LoginCodes[len(r.LoginCodes)-1]
	r.LoginCodes = r.LoginCodes[:len(r.LoginCodes)-1]

	return nil
}
