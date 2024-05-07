package testutil

import (
	"github.com/JosephJoshua/remana-backend/internal/modules/auth/readmodel"
	"github.com/google/uuid"
)

func ModifiedUserDetails(mod func(details *readmodel.UserDetails)) *readmodel.UserDetails {
	details := &readmodel.UserDetails{
		ID:       uuid.New(),
		Username: "not important",
		Role: readmodel.UserDetailsRole{
			ID:           uuid.New(),
			Name:         "not important",
			IsStoreAdmin: true,
		},
		Store: readmodel.UserDetailsStore{
			ID:   uuid.New(),
			Name: "not important",
			Code: "not-important",
		},
	}

	mod(details)
	return details
}
