package notifications

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/grpc"
	"golang.org/x/time/rate"

	pbs "github.com/inhumanLightBackend/notifications/pb"
	grpcg "google.golang.org/grpc"
)

type grpcTransport struct {
	create   grpc.Handler
	findById grpc.Handler
	check    grpc.Handler
}

func NewGRPCServer(ep *Endpoints, logger log.Logger) pbs.NotificationsServer {
	opts := []grpc.ServerOption{
		grpc.ServerErrorLogger(log.With(logger, "tag", "grpc")),
		grpc.ServerBefore(translateMetadataToContext()),
	}

	return &grpcTransport{
		create: grpc.NewServer(
			authenticate()(ep.CreateEndpoint),
			grpcDecodeCreateRequest,
			grpcEncodeCreateResponse,
			opts...,
		),
		findById: grpc.NewServer(
			authenticate()(ep.FindByIdEndpoint),
			grpcDecodeFindByIdRequest,
			grpcEncodeFindByIdResponse,
			opts...,
		),
		check: grpc.NewServer(
			authenticate()(ep.CheckEndpoint),
			grpcDecodeCheckRequest,
			grpcEncodeCheckResponse,
			opts...,
		),
	}
}

func (t *grpcTransport) Create(ctx context.Context, r *pbs.CreateRequest) (*pbs.CreateResponse, error) {
	_, resp, err := t.create.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.CreateResponse), nil
}

func (t *grpcTransport) FindById(ctx context.Context, r *pbs.FindByIdRequest) (*pbs.FindByIdResponse, error) {
	_, resp, err := t.findById.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.FindByIdResponse), nil
}

func (t *grpcTransport) Check(ctx context.Context, r *pbs.CheckRequest) (*pbs.CheckResponse, error) {
	_, resp, err := t.check.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.CheckResponse), nil
}

func NewGRPCClient(conn *grpcg.ClientConn, logger log.Logger) Service {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	opts := []grpc.ClientOption{
		grpc.ClientBefore(translateJwtToMetadata()),
	}

	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = grpc.NewClient(
			conn,
			"pbs.Notifications",
			"Create",
			grpcEncodeCreateRequest,
			grpcDecodeCreateResponse,
			pbs.CreateResponse{},
			opts...,
		).Endpoint()
		createEndpoint = limiter(createEndpoint)
	}
	var findByIdEndpoint endpoint.Endpoint
	{
		findByIdEndpoint = grpc.NewClient(
			conn,
			"pbs.Notifications",
			"FindById",
			grpcEncodeFindByIdRequest,
			grpcDecodeFindByIdResponse,
			pbs.FindByIdResponse{},
			opts...,
		).Endpoint()
		findByIdEndpoint = limiter(findByIdEndpoint)
	}
	var checkEndpoint endpoint.Endpoint
	{
		checkEndpoint = grpc.NewClient(
			conn,
			"pbs.Notifications",
			"Check",
			grpcEncodeCheckRequest,
			grpcDecodeCheckResponse,
			pbs.CheckResponse{},
			opts...,
		).Endpoint()
		checkEndpoint = limiter(checkEndpoint)
	}

	return &Endpoints{
		CreateEndpoint:   createEndpoint,
		FindByIdEndpoint: findByIdEndpoint,
		CheckEndpoint:    checkEndpoint,
	}
}

func grpcDecodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.CreateRequest)

	n := protoNotifToModel(*req.Notification)
	return &CreateRequest{
		Notification: &n,
	}, nil
}

func grpcEncodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*CreateResponse)

	if resp.Err != "" {
		return &pbs.CreateResponse{
			Err: resp.Err,
		}, nil
	}

	return &pbs.CreateResponse{
		Ok:  resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcEncodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*CreateRequest)

	np := modelNotifToProto(*req.Notification)
	return &pbs.CreateRequest{
		Notification: &np,
	}, nil
}

func grpcDecodeCreateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.CreateResponse)

	if resp.Err != "" {
		return &CreateResponse{
			Err: resp.Err,
		}, nil
	}

	return &CreateResponse{
		Ok:  resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcDecodeFindByIdRequest(_ context.Context, r interface{}) (interface{}, error) {
	return &FindByIdRequest{}, nil
}

func grpcEncodeFindByIdResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*FindByIdResponse)

	if resp.Err != "" {
		return &pbs.FindByIdResponse{
			Err: resp.Err,
		}, nil
	}

	var notifications []*pbs.Notification
	{
		for _, noti := range resp.Notifications {
			n := modelNotifToProto(*noti)
			notifications = append(notifications, &n)
		}
	}

	return &pbs.FindByIdResponse{
		Notifications: notifications,
		Err:           resp.Err,
	}, nil
}

func grpcEncodeFindByIdRequest(_ context.Context, r interface{}) (interface{}, error) {
	return &pbs.FindByIdRequest{}, nil
}

func grpcDecodeFindByIdResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.FindByIdResponse)

	if resp.Err != "" {
		return &FindByIdResponse{
			Err: resp.Err,
		}, nil
	}

	var notifications []*Notification
	{
		for _, noti := range resp.Notifications {
			n := protoNotifToModel(*noti)
			notifications = append(notifications, &n)
		}
	}

	return &FindByIdResponse{
		Notifications: notifications,
		Err:           resp.Err,
	}, nil
}

func grpcDecodeCheckRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.CheckRequest)

	var indexes []int
	{
		for _, i := range req.Indexes {
			indexes = append(indexes, int(i))
		}
	}

	return &CheckRequest{
		Indexes: indexes,
	}, nil
}

func grpcEncodeCheckResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*CheckResponse)

	if resp.Err != "" {
		return &pbs.CheckResponse{
			Err: resp.Err,
		}, nil
	}

	return &pbs.CheckResponse{
		Ok:  resp.Ok,
		Err: resp.Err,
	}, nil
}

func grpcEncodeCheckRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*CheckRequest)

	var indexes []int32
	{
		for _, i := range req.Indexes {
			indexes = append(indexes, int32(i))
		}
	}

	return &pbs.CheckRequest{
		Indexes: indexes,
	}, nil
}

func grpcDecodeCheckResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.CheckResponse)

	if resp.Err != "" {
		return &CheckResponse{
			Err: resp.Err,
		}, nil
	}

	return &CheckResponse{
		Ok:  resp.Ok,
		Err: resp.Err,
	}, nil
}

func protoNotifToModel(pb pbs.Notification) Notification {
	return Notification{
		ID:        int(pb.Id),
		Message:   pb.Message,
		Checked:   pb.Checked,
		CreatedAt: pb.CreatedAt,
		For:       int(pb.For),
		Status:    pb.Status,
	}
}

func modelNotifToProto(n Notification) pbs.Notification {
	return pbs.Notification{
		Id:        int32(n.ID),
		Checked:   n.Checked,
		CreatedAt: n.CreatedAt,
		For:       int32(n.For),
		Message:   n.Message,
		Status:    n.Status,
	}
}
