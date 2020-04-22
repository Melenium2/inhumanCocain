package gates

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	auth "github.com/inhumanLightBackend/auth/logic"
	cocain "github.com/inhumanLightBackend/cocain"
)

type (
	SingInResponse struct {
		Token     string    `json:"token"`
		Tokentype string    `json:"token_type"`
		CreatedAt time.Time `json:"created_at"`
	}
)

type Gates interface {
	SignIn(context.Context, string, string) (*SingInResponse, error)
	Logout(context.Context, string) error
	Translate(context.Context, string) (string, error)
}

type GateService struct {
	repo   OpaqueRepository
	logger log.Logger
}

func NewGateService(repo OpaqueRepository, logger log.Logger) Gates {
	return &GateService{
		repo:   repo,
		logger: logger,
	}
}

// Authenticate user login and passwod, then generate opaque token with user params
// and return response to user. Opaque token save to db
func (gs *GateService) SignIn(ctx context.Context, login string, password string) (*SingInResponse, error) {
	us, _, err := cocain.UserGRPCService(cocain.UserPort, gs.logger)
	if err != nil {
		return nil, err
	}
	user, err := us.Authenticate(ctx, login, password)
	if err != nil {
		return nil, err
	}
	as, _, err := cocain.AuthGRPCService(cocain.AuthPort, gs.logger)
	if err != nil {
		return nil, err
	}
	jwt, err := as.Generate(ctx, &auth.GenerateRequest{
		UserId: int32(user.ID),
		Role:   user.Role,
	})
	if err != nil {
		return nil, err
	}
	opaque, err := gs.repo.SaveToken(ctx, user.ID, jwt.Token)
	if err != nil {
		return nil, err
	}

	return &SingInResponse{
		Token: opaque.Opaque,
		Tokentype: "opaque",
		CreatedAt: opaque.CreatedAt,
	}, nil
}

// Remove jwt token from db. Remove user access.
func (gs *GateService) Logout(ctx context.Context, opaque string) error {
	return gs.repo.RemoveToken(ctx, opaque)
}

func (gs *GateService) Translate(ctx context.Context, token string) (string, error) {
	opaque, err := gs.repo.FindToken(ctx, token)
	if err != nil {
		return "", err
	}

	return opaque.Jwt, nil 
}
