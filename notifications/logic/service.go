package notifications

import (
	"context"

	"github.com/go-kit/kit/log"
	auth "github.com/inhumanLightBackend/auth/logic"
)

type Service interface {
	Create(context.Context, *Notification) (string, error)
	FindById(context.Context) ([]*Notification, error)
	Check(context.Context, []int) (string, error)
}

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return &service {
		repository:  repo,
		logger: logger,
	}
}

func (s *service) Create(ctx context.Context, n *Notification) (string, error) {
	if err := s.repository.CreateNotification(ctx, n); err != nil {
		return "", err
	}

	return "Created", nil
}

func (s *service) FindById(ctx context.Context) ([]*Notification, error) {
	claims, err := getClaims(ctx, s.logger)
	if err != nil {
		return nil, err
	}

	return s.repository.FindNotificationsById(ctx, int(claims.UserId))
}

func (s *service) Check(ctx context.Context, indexes []int) (string, error) {
	claims, err := getClaims(ctx, s.logger)
	if err != nil {
		return "", err
	}

	if err := s.repository.CheckNotification(ctx, indexes, int(claims.UserId)); err != nil {
		return "", err
	}

	return "Checked", nil
}

func getClaims(ctx context.Context, logger log.Logger) (*auth.ValidateResponse, error) {
	var claims *auth.ValidateResponse 
	{
		userRawId := ctx.Value(auth.CtxUserKey).(string)
		serv, _, err := AuthGRPCService(AuthserverPort, logger)
		if err != nil {
			return nil, err
		}
		claims, err = serv.Validate(ctx, &auth.ValidateRequest{
			Token: userRawId,
		})
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}
