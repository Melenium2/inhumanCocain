package auth

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/grpc"
	pbs "github.com/inhumanLightBackend/auth/pb"
	"golang.org/x/time/rate"
	grpcg "google.golang.org/grpc"
)

type grpcTransport struct {
	generate grpc.Handler
	validate grpc.Handler
}

func NewGRPCServer(ep *Endpoints, logger log.Logger) pbs.AuthServer {
	opts := []grpc.ServerOption {
		grpc.ServerErrorLogger(log.With(logger, "tag", "grpc")),
	}

	return &grpcTransport {
		generate: grpc.NewServer(
			ep.GenerateEndpoint,
			grpcDecodeGenerateRequest,
			grpcEncodeGenerateResponse,
			opts...,
		),
		validate: grpc.NewServer(
			ep.ValidateEndpoint,
			grpcDecodeValidateRequest,
			grpcEncodeValidateResponse,
			opts...,
		),
	}
}

func (t *grpcTransport) Generate(ctx context.Context, r *pbs.GenerateRequest) (*pbs.GenerateResponse, error) {
	_, resp, err := t.generate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.GenerateResponse), nil
}

func (t *grpcTransport) Validate(ctx context.Context, r *pbs.ValidateRequest) (*pbs.ValidateResponse, error) {
	_, resp, err := t.validate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.ValidateResponse), nil
}

func NewGRPCClient(conn *grpcg.ClientConn, logger log.Logger) Service {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	var generateEndpoint endpoint.Endpoint
	{
		generateEndpoint = grpc.NewClient(
			conn,
			"pbs.Auth",
			"Generate",
			grpcEncodeGenerateRequest,
			grpcDecodeGenerateResponse,
			pbs.GenerateResponse{},
		).Endpoint()
		generateEndpoint = limiter(generateEndpoint)
	}
	var validateEndpoint endpoint.Endpoint
	{
		validateEndpoint = grpc.NewClient(
			conn,
			"pbs.Auth",
			"Validate",
			grpcEncodeValidateRequest,
			grpcDecodeValidateResponse,
			pbs.ValidateResponse{},
		).Endpoint()
		validateEndpoint = limiter(validateEndpoint)
	}

	return &Endpoints{
		GenerateEndpoint: generateEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
}

func grpcDecodeGenerateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.GenerateRequest)

	return &GenerateRequest {
		UserId: req.UserId,
		Role: req.Role,
	}, nil
}

func grpcEncodeGenerateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*GenerateResponse)

	if resp.Err != "" {
		println(resp.Err)
		return &pbs.GenerateResponse {
			Token: "error",
		}, nil
	}

	return &pbs.GenerateResponse {
		Token: resp.Token,
	}, nil
}

func grpcEncodeGenerateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*GenerateRequest)
	
	return &pbs.GenerateRequest {
		UserId: req.UserId,
		Role: req.Role,
	}, nil
}

func grpcDecodeGenerateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.GenerateResponse)
	
	return &GenerateResponse {
		Token: resp.Token,
	}, nil
}

func grpcDecodeValidateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.ValidateRequest)

	return &ValidateRequest {
		Token: req.Token,
	}, nil
}

func grpcEncodeValidateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*ValidateResponse)

	if resp.Err != "" {
		println(resp.Err)
		return &pbs.ValidateResponse {
			Err: resp.Err,
		}, nil
	}

	return &pbs.ValidateResponse {
		Role: resp.Role,
		Token: resp.Token,
		UserId: resp.UserId,
	}, nil
}

func grpcEncodeValidateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*ValidateRequest)

	return &pbs.ValidateRequest {
		Token: req.Token,
	}, nil
}

func grpcDecodeValidateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.ValidateResponse)

	if resp.Err != "" {
		return &ValidateResponse {
			Err: resp.Err,
		}, nil
	}

	return &ValidateResponse {
		UserId: resp.UserId,
		Role: resp.Role,
		Token: resp.Token,
	}, nil
}