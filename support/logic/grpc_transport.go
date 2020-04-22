package support

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/grpc"
	"golang.org/x/time/rate"

	pbs "github.com/inhumanLightBackend/support/pb"
	grpcg "google.golang.org/grpc"
)

type grpcTransport struct {
	createTicket grpc.Handler
	getTicket    grpc.Handler
	tickets      grpc.Handler
	acceptTicket grpc.Handler
	addMessage   grpc.Handler
	getMessages  grpc.Handler
	changeStatus grpc.Handler
}

func NewGRPCServer(ep *Endpoints, logger log.Logger) pbs.SupportServer {
	opts := []grpc.ServerOption{
		grpc.ServerErrorLogger(log.With(logger, "tag", "grpc")),
		grpc.ServerBefore(translateMetadataToContext()),
	}

	return &grpcTransport {
		createTicket: grpc.NewServer(
			authenticate()(ep.CreateTicketEndpoint),
			grpcDecodeCreateTicketRequest,
			grpcEncodeCreateTicketResponse,
			opts...,
		),
		getTicket: grpc.NewServer(
			authenticate()(ep.GetTicketEndpoint),
			grpcDecodeGetTicketRequest,
			grpcEncodeGetTicketResponse,
			opts...,
		),
		tickets: grpc.NewServer(
			authenticate()(ep.TicketsEndpoint),
			grpcDecodeTicketsRequest,
			grpcEncodeTicketsResponse,
			opts...,
		),
		acceptTicket: grpc.NewServer(
			authenticate()(ep.AcceptTicketEndpoint),
			grpcDecodeAcceptTicketRequest,
			grpcEncodeAcceptTicketResponse,
			opts...,
		),
		addMessage: grpc.NewServer(
			authenticate()(ep.AddMessageEndpoint),
			grpcDecodeAddMessageRequest,
			grpcEncodeAddMessageResponse,
			opts...,
		),
		getMessages: grpc.NewServer(
			authenticate()(ep.GetMessagesEndpoint),
			grpcDecodeGetMessagesRequest,
			grpcEncodeGetMessagesResponse,
			opts...,
		),
		changeStatus: grpc.NewServer(
			authenticate()(ep.ChangeStatusEndpoint),
			grpcDecodeChangeStatusRequest,
			grpcEncodeChangeStatusResponse,
			opts...,
		),
	}
}

func (t *grpcTransport) CreateTicket(ctx context.Context, r *pbs.CreateTicketRequest) (*pbs.CreateTicketResponse, error) {
	_, resp, err := t.createTicket.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.CreateTicketResponse), nil
}

func (t *grpcTransport) GetTicket(ctx context.Context, r *pbs.GetTicketRequest) (*pbs.GetTicketResponse, error) {
	_, resp, err := t.getTicket.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.GetTicketResponse), nil
}

func (t *grpcTransport) Tickets(ctx context.Context, r *pbs.TicketsRequest) (*pbs.TicketsResponse, error) {
	_, resp, err := t.tickets.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.TicketsResponse), nil
}

func (t *grpcTransport) AcceptTicket(ctx context.Context, r *pbs.AcceptTicketRequest) (*pbs.AcceptTicketResponse, error) {
	_, resp, err := t.acceptTicket.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.AcceptTicketResponse), nil
}

func (t *grpcTransport) AddMessage(ctx context.Context, r *pbs.AddMessageRequest) (*pbs.AddMessageResponse, error) {
	_, resp, err := t.addMessage.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.AddMessageResponse), nil
}

func (t *grpcTransport) GetMessages(ctx context.Context, r *pbs.GetMessagesRequest) (*pbs.GetMessagesResponse, error) {
	_, resp, err := t.getMessages.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.GetMessagesResponse), nil
}

func (t *grpcTransport) ChangeStatus(ctx context.Context, r *pbs.ChangeStatusRequest) (*pbs.ChangeStatusResponse, error) {
	_, resp, err := t.changeStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.ChangeStatusResponse), nil
}

func NewGRPCClient(conn *grpcg.ClientConn, logger log.Logger) Service {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	opts := []grpc.ClientOption{
		grpc.ClientBefore(translateJwtToMetadata()),
	}
	
	var createTicketEndpoint endpoint.Endpoint
	{
		createTicketEndpoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"CreateTicket",
			grpcEncodeCreateTicketRequest,
			grpcDecodeCreateTicketResponse,
			pbs.CreateTicketResponse {},
			opts...,
		).Endpoint()
		createTicketEndpoint = limiter(createTicketEndpoint)
	}
	var getTicketEndpoint endpoint.Endpoint
	{
		getTicketEndpoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"GetTicket",
			grpcEncodeGetTicketRequest,
			grpcDecodeGetTicketResponse,
			pbs.GetTicketResponse {},
			opts...,
		).Endpoint()
		getTicketEndpoint = limiter(getTicketEndpoint)
	}
	var ticketsEndoint endpoint.Endpoint 
	{
		ticketsEndoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"Tickets",
			grpcEncodeTicketsRequest,
			grpcDecodeTicketsResponse,
			pbs.TicketsResponse {},
			opts...,
		).Endpoint()
		ticketsEndoint = limiter(ticketsEndoint)
	}
	var acceptTicketEndpoint endpoint.Endpoint 
	{
		acceptTicketEndpoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"AcceptTicket",
			grpcEncodeAcceptTicketRequest,
			grpcDecodeAcceptTicketResponse,
			pbs.AcceptTicketResponse {},
			opts...,
		).Endpoint()
		acceptTicketEndpoint = limiter(acceptTicketEndpoint)
	}
	var addMessageEndpoint endpoint.Endpoint
	{
		addMessageEndpoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"AddMessage",
			grpcEncodeAddMessageRequest,
			grpcDecodeAddMessageResponse,
			pbs.AddMessageResponse {},
			opts...,
		).Endpoint()
		addMessageEndpoint = limiter(addMessageEndpoint)
	}
	var getMessagesEndpoint endpoint.Endpoint
	{
		getMessagesEndpoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"GetMessages",
			grpcEncodeGetMessagesRequest,
			grpcDecodeGetMessagesResponse,
			pbs.GetMessagesResponse {},
			opts...,
		).Endpoint()
		getMessagesEndpoint = limiter(getMessagesEndpoint)
	}
	var changeStatusEndpoint endpoint.Endpoint
	{
		changeStatusEndpoint = grpc.NewClient(
			conn,
			"pbs.Support",
			"ChangeStatus",
			grpcEncodeChangeStatusRequest,
			grpcDecodeChangeStatusResponse,
			pbs.ChangeStatusResponse {},
			opts...,
		).Endpoint()
		changeStatusEndpoint = limiter(changeStatusEndpoint)
	}

	return &Endpoints{
		CreateTicketEndpoint: createTicketEndpoint,
		GetTicketEndpoint: getTicketEndpoint,
		TicketsEndpoint: ticketsEndoint,
		AcceptTicketEndpoint: acceptTicketEndpoint,
		AddMessageEndpoint: addMessageEndpoint,
		GetMessagesEndpoint: getMessagesEndpoint,
		ChangeStatusEndpoint: changeStatusEndpoint,
	}
}

func grpcDecodeCreateTicketRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.CreateTicketRequest)
	
	t := protoTicketToModel(*req.Ticket)
	return &CreateTicketRequest {
		Ticket: &t,
	}, nil
}

func grpcEncodeCreateTicketRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*CreateTicketRequest)

	pt := modelTicketToProto(*req.Ticket)
	return &pbs.CreateTicketRequest {
		Ticket: &pt,
	}, nil
}

func grpcDecodeCreateTicketResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.CreateTicketResponse)

	if resp.Err != "" {
		return &CreateTicketResponse {
			Err: resp.Err,
		}, nil
	}

	println("grpcDecodeCreateTicketResponse", resp.Ok, resp.Err)
	return &CreateTicketResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcEncodeCreateTicketResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*CreateTicketResponse)

	if resp.Err != "" {
		return &pbs.CreateTicketResponse {
			Err: resp.Err,
		}, nil
	}
	
	println("grpcEncodeCreateTicketResponse", resp.Ok, resp.Err)
	return &pbs.CreateTicketResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcDecodeGetTicketRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.GetTicketRequest)
	
	return &GetTicketRequest {
		TicketId: int(req.TicketId),
	}, nil
}

func grpcEncodeGetTicketRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*GetTicketRequest)
	
	return &pbs.GetTicketRequest {
		TicketId: int32(req.TicketId),
	}, nil
}

func grpcDecodeGetTicketResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.GetTicketResponse)
	
	if resp.Err != "" {
		return &GetTicketResponse {
			Err: resp.Err,
		}, nil
	}
	
	t := protoTicketToModel(*resp.Ticket)
	return &GetTicketResponse {
		Ticket: &t,
		Err: resp.Err,
	}, nil
}

func grpcEncodeGetTicketResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*GetTicketResponse)
	
	if resp.Err != "" {
		return &pbs.GetTicketResponse {
			Err: resp.Err,
		}, nil
	}
	
	pt := modelTicketToProto(*resp.Ticket)
	return &pbs.GetTicketResponse {
		Ticket: &pt,
		Err: resp.Err,
	}, nil
}

func grpcDecodeTicketsRequest(_ context.Context, r interface{}) (interface{}, error) {
	return &TicketsRequest {}, nil
}

func grpcEncodeTicketsRequest(_ context.Context, r interface{}) (interface{}, error) {
	return &pbs.TicketsRequest {}, nil
}

func grpcDecodeTicketsResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.TicketsResponse)
	
	if resp.Err != "" {
		return &TicketsResponse {
			Err: resp.Err,
		}, nil
	}
	
	var tickets []*Ticket
	{
		for _, ticket := range resp.Tickets {
			pt := protoTicketToModel(*ticket)
			tickets = append(tickets, &pt)
		}
	}

	return &TicketsResponse {
		Tickets: tickets,
		Err: resp.Err,
	}, nil
}

func grpcEncodeTicketsResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*TicketsResponse)
	
	if resp.Err != "" {
		return &pbs.TicketsResponse {
			Err: resp.Err,
		}, nil
	}
	
	var tickets []*pbs.Ticket
	{
		for _, ticket := range resp.Tickets {
			t := modelTicketToProto(*ticket)
			tickets = append(tickets, &t)
		}
	}

	return &pbs.TicketsResponse {
		Tickets: tickets,
		Err: resp.Err,
	}, nil
}

func grpcDecodeAcceptTicketRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.AcceptTicketRequest)
	
	return &AcceptTicketRequest {
		TicketId: int(req.TicketId),
	}, nil
}

func grpcEncodeAcceptTicketRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*AcceptTicketRequest)
	
	return &pbs.AcceptTicketRequest {
		TicketId: int32(req.TicketId),
	}, nil
}

func grpcDecodeAcceptTicketResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.AcceptTicketResponse)
	
	if resp.Err != "" {
		return &AcceptTicketResponse {
			Err: resp.Err,
		}, nil
	}
	
	return &AcceptTicketResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcEncodeAcceptTicketResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*AcceptTicketResponse)
	
	if resp.Err != "" {
		return &pbs.AcceptTicketResponse {
			Err: resp.Err,
		}, nil
	}
	
	return &pbs.AcceptTicketResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcDecodeAddMessageRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.AddMessageRequest)
	
	m := protoMessageToModel(*req.Message)
	return &AddMessageRequest {
		Message: &m,
	}, nil
}

func grpcEncodeAddMessageRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*AddMessageRequest)
	
	pm := modelMessageToProto(*req.Message)
	return &pbs.AddMessageRequest {
		Message: &pm,
	}, nil
}

func grpcDecodeAddMessageResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.AddMessageResponse)
	
	if resp.Err != "" {
		return &AddMessageResponse {
			Err: resp.Err,
		}, nil
	}
	
	return &AddMessageResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcEncodeAddMessageResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*AddMessageResponse)
	
	if resp.Err != "" {
		return &pbs.AddMessageResponse {
			Err: resp.Err,
		}, nil
	}
	
	return &pbs.AddMessageResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcDecodeGetMessagesRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.GetMessagesRequest)
	
	return &GetMessagesRequest {
		TicketId: int(req.TicketId),
	}, nil
}

func grpcEncodeGetMessagesRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*GetMessagesRequest)
	
	return &pbs.GetMessagesRequest {
		TicketId: int32(req.TicketId),
	}, nil
}

func grpcDecodeGetMessagesResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.GetMessagesResponse)
	
	if resp.Err != "" {
		return &GetMessagesResponse {
			Err: resp.Err,
		}, nil
	}
	
	var messages []*TicketMessage
	{
		for _, message := range resp.Messages {
			m := protoMessageToModel(*message)
			messages = append(messages, &m)
		}
	}

	return &GetMessagesResponse {
		Messages: messages,
		Err: resp.Err,
	}, nil
}

func grpcEncodeGetMessagesResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*GetMessagesResponse)
	
	if resp.Err != "" {
		return &pbs.GetMessagesResponse {
			Err: resp.Err,
		}, nil
	}

	var messages []*pbs.TicketMessage
	{
		for _, message := range resp.Messages {
			pm := modelMessageToProto(*message)
			messages = append(messages, &pm)
		}
	}
	
	return &pbs.GetMessagesResponse {
		Messages: messages,
		Err: resp.Err,
	}, nil
}

func grpcDecodeChangeStatusRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.ChangeStatusRequest)
	
	return &ChangeStatusRequest {
		TicketId: int(req.TicketId),
		Status: req.Status,
	}, nil
}

func grpcEncodeChangeStatusRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*ChangeStatusRequest)
	
	return &pbs.ChangeStatusRequest {
		TicketId: int32(req.TicketId),
		Status: req.Status,
	}, nil
}

func grpcDecodeChangeStatusResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.ChangeStatusResponse)
	
	if resp.Err != "" {
		return &ChangeStatusResponse {
			Err: resp.Err,
		}, nil
	}
	
	return &ChangeStatusResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcEncodeChangeStatusResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*ChangeStatusResponse)
	
	if resp.Err != "" {
		return &pbs.ChangeStatusResponse {
			Err: resp.Err,
		}, nil
	}
	
	return &pbs.ChangeStatusResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func modelTicketToProto(t Ticket) pbs.Ticket {
	return pbs.Ticket{
		Id: int32(t.ID),
		CreatedAt: t.CreatedAt,
		Description: t.Description,
		From: int32(t.From),
		Helper: int32(t.Helper),
		Section: t.Section,
		Status: t.Status,
		Title: t.Title,
	}
}

func protoTicketToModel(pt pbs.Ticket) Ticket {
	return Ticket{
		ID: int(pt.Id),
		CreatedAt: pt.CreatedAt,
		Description: pt.Description,	
		From: int(pt.From),
		Helper: int(pt.Helper),
		Section: pt.Section,
		Status: pt.Status,
		Title: pt.Title,
	}
}

func modelMessageToProto(m TicketMessage) pbs.TicketMessage {
	return pbs.TicketMessage{
		Id: int32(m.ID),
		Message: m.Message,
		SendedAt: m.SendedAt,
		TicketId: int32(m.TicketId),
		Who: int32(m.Who),
	}
}

func protoMessageToModel(pm pbs.TicketMessage) TicketMessage {
	return TicketMessage{
		ID: int(pm.Id),
		Message: pm.Message,
		SendedAt: pm.SendedAt,
		TicketId: int(pm.TicketId),
		Who: int(pm.Who),
	}
}