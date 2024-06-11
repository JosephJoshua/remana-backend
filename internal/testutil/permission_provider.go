package testutil

import (
	"context"

	"github.com/JosephJoshua/remana-backend/internal/modules/permission"
	"github.com/google/uuid"
)

type PermissionProviderStub struct {
	roleID          uuid.UUID
	rolePermissions []permission.Permission
	err             error
}

func NewPermissionProviderStub(
	roleID uuid.UUID,
	rolePermissions []permission.Permission,
	err error,
) *PermissionProviderStub {
	return &PermissionProviderStub{
		roleID:          roleID,
		rolePermissions: rolePermissions,
		err:             err,
	}
}

func (p *PermissionProviderStub) SetError(err error) {
	p.err = err
}

func (p *PermissionProviderStub) Can(
	_ context.Context,
	roleID uuid.UUID,
	permission permission.Permission,
) (bool, error) {
	if p.err != nil {
		return false, p.err
	}

	if p.roleID != roleID {
		return false, nil
	}

	for _, p := range p.rolePermissions {
		if p.GroupName() == permission.GroupName() && p.Name() == permission.Name() {
			return true, nil
		}
	}

	return false, nil
}
