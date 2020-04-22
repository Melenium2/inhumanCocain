package support

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateTicketEndpoint endpoint.Endpoint
	GetTicketEndpoint    endpoint.Endpoint
	TicketsEndpoint      endpoint.Endpoint
	AcceptTicketEndpoint endpoint.Endpoint
	AddMessageEndpoint   endpoint.Endpoint
	GetMessagesEndpoint  endpoint.Endpoint
	ChangeStatusEndpoint endpoint.Endpoint
}

type (
	CreateTicketRequest struct {
		Ticket *Ticket `json:"ticket"`
	}

	CreateTicketResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}

	GetTicketRequest struct {
		TicketId int `json:"ticket_id"`
	}

	GetTicketResponse struct {
		Ticket *Ticket `json:"ticket,omitempty"`
		Err    string  `json:"error,omitempty"`
	}

	TicketsRequest struct {
	}

	TicketsResponse struct {
		Tickets []*Ticket `json:"tickets,omitempty"`
		Err     string    `json:"error,omitempty"`
	}

	AcceptTicketRequest struct {
		TicketId int `json:"ticket_id"`
	}

	AcceptTicketResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}

	AddMessageRequest struct {
		Message *TicketMessage `json:"message"`
	}

	AddMessageResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}

	GetMessagesRequest struct {
		TicketId int `json:"ticket_id"`
	}

	GetMessagesResponse struct {
		Messages []*TicketMessage `json:"messages,omitempty"`
		Err      string           `json:"error,omitempty"`
	}

	ChangeStatusRequest struct {
		TicketId int    `json:"ticket_id"`
		Status   string `json:"status"`
	}

	ChangeStatusResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}
)

func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		CreateTicketEndpoint: MakeCreateTicketEndpoint(s),
		GetTicketEndpoint: MakeGetTicketEndpoint(s),
		TicketsEndpoint: MakeTicketsEndpoint(s),
		AcceptTicketEndpoint: MakeAcceptTicketEndpoint(s),
		AddMessageEndpoint: MakeAddMessageEndpoint(s),
		GetMessagesEndpoint: MakeGetMessagesEndpoint(s),
		ChangeStatusEndpoint: MakeChangeStatusEndpoint(s),
	}
}

func (ep *Endpoints) CreateTicket(ctx context.Context, t *Ticket) (string, error) {
	r, err := ep.CreateTicketEndpoint(ctx, &CreateTicketRequest{
		Ticket: t,
	})
	if err != nil {
		return "", err
	}
	resp := r.(*CreateTicketResponse)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.Ok, nil
}

func (ep *Endpoints) GetTicket(ctx context.Context, ticketId int) (*Ticket, error) {
	r, err := ep.GetTicketEndpoint(ctx, &GetTicketRequest{
		TicketId: ticketId,
	})
	if err != nil {
		return nil, err
	}
	resp := r.(*GetTicketResponse)
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}

	return resp.Ticket, nil
}

func (ep *Endpoints) Tickets(ctx context.Context) ([]*Ticket, error) {
	r, err := ep.TicketsEndpoint(ctx, &TicketsRequest{})
	if err != nil {
		return nil, err
	}
	resp := r.(*TicketsResponse)
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}
	return resp.Tickets, nil
}

func (ep *Endpoints) AcceptTicket(ctx context.Context, ticketId int) (string, error) {
	r, err := ep.AcceptTicketEndpoint(ctx, &AcceptTicketRequest{
		TicketId: ticketId,
	})
	if err != nil {
		return "", err
	}
	resp := r.(*AcceptTicketResponse)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.Ok, nil
}

func (ep *Endpoints) AddMessage(ctx context.Context, m *TicketMessage) (string, error) {
	r, err := ep.AddMessageEndpoint(ctx, &AddMessageRequest{
		Message: m,
	})
	if err != nil {
		return "", err
	}
	resp := r.(*AddMessageResponse)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.Ok, nil
}

func (ep *Endpoints) GetMessages(ctx context.Context, ticketId int) ([]*TicketMessage, error) {
	r, err := ep.GetMessagesEndpoint(ctx, &GetMessagesRequest{
		TicketId: ticketId,
	})
	if err != nil {
		return nil, err
	}
	resp := r.(*GetMessagesResponse)
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}
	return resp.Messages, nil
}

func (ep *Endpoints) ChangeStatus(ctx context.Context, ticketId int, status string) (string, error) {
	r, err := ep.ChangeStatusEndpoint(ctx, &ChangeStatusRequest{
		TicketId: ticketId,
		Status: status,
	})
	if err != nil {
		return "", err
	}
	resp := r.(*ChangeStatusResponse)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.Ok, nil
}

func MakeCreateTicketEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CreateTicketRequest)
		res, err := s.CreateTicket(ctx, req.Ticket)
		if err != nil {
			return &CreateTicketResponse {
				Err: err.Error(),
			}, nil
		}

		return &CreateTicketResponse {
			Ok: res,
		}, nil
	}
}

func MakeGetTicketEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetTicketRequest)
		res, err := s.GetTicket(ctx, req.TicketId)
		if err != nil {
			return &GetTicketResponse {
				Err: err.Error(),
			}, nil
		}

		return &GetTicketResponse {
			Ticket: res,
		}, nil
	}
}

func MakeTicketsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.Tickets(ctx)
		if err != nil {
			return &TicketsResponse {
				Err: err.Error(),
			}, nil
		}

		return &TicketsResponse {
			Tickets: res,
		}, nil
	}
}

func MakeAcceptTicketEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*AcceptTicketRequest)
		res, err := s.AcceptTicket(ctx, req.TicketId)
		if err != nil {
			return &AcceptTicketResponse {
				Err: err.Error(),
			}, nil
		}

		return &AcceptTicketResponse {
			Ok: res,
		}, nil
	}
}

func MakeAddMessageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*AddMessageRequest)
		res, err := s.AddMessage(ctx, req.Message)
		if err != nil {
			return &AddMessageResponse {
				Err: err.Error(),
			}, nil
		}

		return &AddMessageResponse {
			Ok: res,
		}, nil
	}
}

func MakeGetMessagesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GetMessagesRequest)
		res, err := s.GetMessages(ctx, req.TicketId)
		if err != nil {
			return &GetMessagesResponse {
				Err: err.Error(),
			}, nil
		}

		return &GetMessagesResponse {
			Messages: res,
		}, nil
	}
}

func MakeChangeStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*ChangeStatusRequest)
		res, err := s.ChangeStatus(ctx, req.TicketId, req.Status)
		if err != nil {
			return &ChangeStatusResponse {
				Err: err.Error(),
			}, nil
		}

		return &ChangeStatusResponse {
			Ok: res,
		}, nil
	}	
}