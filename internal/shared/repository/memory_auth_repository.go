package repository

import (
	"context"
	"slices"

	"github.com/JosephJoshua/repair-management-backend/internal/shared/apperror"
	"github.com/JosephJoshua/repair-management-backend/internal/shared/domain"
)

type MemoryAuthRepository struct {
	users []domain.User
}

func NewMemoryAuthRepository(users []domain.User) *MemoryAuthRepository {
	return &MemoryAuthRepository{users: users}
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
		return domain.User{}, apperror.ErrUserNotFound
	}

	return r.users[i], nil
}
