package support

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	auth "github.com/inhumanLightBackend/auth/logic"
	user "github.com/inhumanLightBackend/user/logic"
)

type Service interface {
	CreateTicket(context.Context, *Ticket) (string, error)
	GetTicket(context.Context, int) (*Ticket, error)
	Tickets(context.Context) ([]*Ticket, error)
	AcceptTicket(context.Context, int) (string, error)
	AddMessage(context.Context, *TicketMessage) (string, error)
	GetMessages(context.Context, int) ([]*TicketMessage, error)
	ChangeStatus(context.Context, int, string) (string, error)
}

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(r Repository, l log.Logger) Service {
	return &service{
		repository: r,
		logger:     l,
	}
}

func (s *service) CreateTicket(ctx context.Context, t *Ticket) (string, error) {
	claims, err := getClaims(ctx, s.logger)
	if err != nil {
		return "", err
	}
	t.From = int(claims.UserId)

	if err := s.repository.Create(ctx, t); err != nil {
		return "", err
	}

	return "Created", nil
}

func (s *service) GetTicket(ctx context.Context, ticketId int) (*Ticket, error) {
	return s.repository.Find(ctx, ticketId)
}

func (s *service) Tickets(ctx context.Context) ([]*Ticket, error) {
	claims, err := getClaims(ctx, s.logger)
	if err != nil {
		return nil, err
	}

	return s.repository.FindAll(ctx, int(claims.UserId))
}

func (s *service) AcceptTicket(ctx context.Context, ticketId int) (string, error) {
	claims, err := getClaims(ctx, s.logger)
	if err != nil {
		return "", err
	}
	if claims.Role != user.ADMIN {
		return "", errPermissionsDeni
	}

	if err := s.repository.Accept(ctx, ticketId, int(claims.UserId)); err != nil {
		return "", err
	}

	return "Accepted", nil
}

func (s *service) AddMessage(ctx context.Context, m *TicketMessage) (string, error) {
	claims, err := getClaims(ctx, s.logger)
	if err != nil {
		return "", err
	}
	m.Who = int(claims.UserId)

	_, err = s.repository.Find(ctx, m.TicketId)
	if err != nil {
		return "", err
	}

	if err := s.repository.AddMessage(ctx, m); err != nil {
		return "", nil
	}

	return "Added", nil
}

func (s *service) GetMessages(ctx context.Context, ticketId int) ([]*TicketMessage, error) {
	return s.repository.TakeMessages(ctx, ticketId)
}

func (s *service) ChangeStatus(ctx context.Context, ticketId int, status string) (string, error) {
	switch status {
	case Opened:
	case InProcess:
	case Closed:
	default:
		return "", errBadRequest
	}

	if err := s.repository.ChangeStatus(ctx, ticketId, status); err != nil {
		return "", nil
	}

	return fmt.Sprintf("Status changed to %s", status), nil
}

func getClaims(ctx context.Context, l log.Logger) (*auth.ValidateResponse, error) {
	var claims *auth.ValidateResponse
	{
		userRawClaims := ctx.Value(auth.CtxUserKey).(string)
		serv, _, err := AuthGRPCService(AuthserverPort, l)
		if err != nil {
			return nil, err
		}
		claims, err = serv.Validate(ctx, &auth.ValidateRequest{
			Token: userRawClaims,
		})
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}
