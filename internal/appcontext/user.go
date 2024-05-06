package appcontext

import (
	"context"

	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
)

type userCtxKey struct{}

func NewContextWithUser(ctx context.Context, user *readmodel.UserDetails) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func GetUserFromContext(ctx context.Context) (*readmodel.UserDetails, bool) {
	user, ok := ctx.Value(userCtxKey{}).(*readmodel.UserDetails)
	return user, ok
}
