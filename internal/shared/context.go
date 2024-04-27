package shared

import (
	"context"
	"net/url"

	"github.com/JosephJoshua/remana-backend/internal/shared/readmodel"
)

type userCtxKey struct{}

func NewContextWithUser(ctx context.Context, user *readmodel.UserDetails) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func GetUserFromContext(ctx context.Context) (*readmodel.UserDetails, bool) {
	user, ok := ctx.Value(userCtxKey{}).(*readmodel.UserDetails)
	return user, ok
}

type requestURLCtxKey struct{}

func NewContextWithRequestURL(ctx context.Context, url url.URL) context.Context {
	return context.WithValue(ctx, requestURLCtxKey{}, url)
}

func GetRequestURLFromContext(ctx context.Context) (url.URL, bool) {
	url, ok := ctx.Value(requestURLCtxKey{}).(url.URL)
	return url, ok
}
