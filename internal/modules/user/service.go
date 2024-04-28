package user

import (
	"context"
	"net/http"

	"github.com/JosephJoshua/remana-backend/internal/apierror"
	"github.com/JosephJoshua/remana-backend/internal/genapi"
	"github.com/JosephJoshua/remana-backend/internal/shared"
	"github.com/rs/zerolog"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetMyUserDetails(ctx context.Context) (*genapi.UserDetails, error) {
	l := zerolog.Ctx(ctx)

	user, ok := shared.GetUserFromContext(ctx)
	if !ok {
		l.Error().Msg("GetUserFromContext(); failed to get user from context")
		return nil, apierror.ToAPIError(http.StatusInternalServerError, "failed to get user from context")
	}

	return &genapi.UserDetails{
		ID:       user.ID,
		Username: user.Username,
		Role: genapi.UserDetailsRole{
			ID:           user.Role.ID,
			Name:         user.Role.Name,
			IsStoreAdmin: user.Role.IsStoreAdmin,
		},
		Store: genapi.UserDetailsStore{
			ID:   user.Store.ID,
			Name: user.Store.Name,
			Code: user.Store.Code,
		},
	}, nil
}
