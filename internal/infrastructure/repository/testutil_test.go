package repository_test

import (
	"context"

	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/google/uuid"
)

type permissionProviderStub struct{}

func (p permissionProviderStub) Can(_ context.Context, _ uuid.UUID, _ permission.Permission) (bool, error) {
	return true, nil
}
