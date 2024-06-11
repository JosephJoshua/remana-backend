package permission

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ProviderRepository interface {
	IsStoreAdmin(ctx context.Context, roleID uuid.UUID) (bool, error)
	HasPermission(ctx context.Context, roleID uuid.UUID, permissionGroupName string, permissionName string) (bool, error)
}

type Provider interface {
	Can(ctx context.Context, roleID uuid.UUID, permission Permission) (bool, error)
}

type provider struct {
	repo ProviderRepository
}

func NewProvider(repo ProviderRepository) Provider {
	return &provider{repo: repo}
}

func (p *provider) Can(ctx context.Context, roleID uuid.UUID, permission Permission) (bool, error) {
	l := zerolog.Ctx(ctx)

	if ok, err := p.repo.IsStoreAdmin(ctx, roleID); err != nil {
		l.Error().Err(err).Msg("failed to check if role is store admin")
		return false, fmt.Errorf("failed to check if role is store admin: %w", err)
	} else if ok {
		return true, nil
	}

	ok, err := p.repo.HasPermission(ctx, roleID, permission.GroupName(), permission.Name())
	if err != nil {
		l.Error().Err(err).Msg("failed to check if role has permission")
		return false, fmt.Errorf("failed to check if role has permission: %w", err)
	}

	return ok, nil
}
