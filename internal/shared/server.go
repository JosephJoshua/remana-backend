package shared

import (
	"context"

	"github.com/JosephJoshua/repair-management-backend/internal/genapi"
)

type Server struct{}

func (s Server) Login(ctx context.Context, req *genapi.LoginCredentials) (genapi.LoginRes, error) {
	return &genapi.LoginNoContent{}, nil
}

func (s Server) LoginCodePrompt(ctx context.Context, req *genapi.LoginCodePrompt, params genapi.LoginCodePromptParams) (*genapi.LoginCodePromptNoContent, error) {
	return &genapi.LoginCodePromptNoContent{}, nil
}

func (s Server) NewError(ctx context.Context, err error) *genapi.ErrorStatusCode {
	return &genapi.ErrorStatusCode{}
}
