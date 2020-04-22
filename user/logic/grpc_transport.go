package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/grpc"
	pbs "github.com/inhumanLightBackend/user/pb"
	"golang.org/x/time/rate"
	grpcg "google.golang.org/grpc"
)

type grpcTransport struct {
	authenticate grpc.Handler
	create      grpc.Handler
	findByEmail grpc.Handler
	findById    grpc.Handler
	update      grpc.Handler
}

func NewGRPCServer(ep *Endpoints, logger log.Logger) pbs.UsersServer {
	opts := []grpc.ServerOption {
		grpc.ServerErrorLogger(log.With(logger, "tag", "grpc")),
	}
	return &grpcTransport {
		authenticate: grpc.NewServer(
			ep.AuthenticateEndpoint,
			grpcDecodeAuthenticateRequest,
			grpcEncodeAuthenticateResponse,
			opts...,
		),
		create: grpc.NewServer(
			ep.CreateEndpoint,
			decodeGRPCCreateUserRequest,
			encodeGRPCCreateUserResponse,
		),
		findByEmail: grpc.NewServer(
			authenticate()(ep.FindByEmailEndpoint),
			decodeGRPCFindByEmailRequest,
			encodeGRPCFindByEmailResponse,
			append(opts, grpc.ServerBefore(translateMetadataToContext()))...,
		),
		findById: grpc.NewServer(
			authenticate()(ep.FindByIdEndpoint),
			decodeGRPCFindByIdRequest,
			encodeGRPCFindByIdResponse,
			append(opts, grpc.ServerBefore(translateMetadataToContext()))...,
		),
		update: grpc.NewServer(
			authenticate()(ep.UpdateEndpoint),
			decodeGRPCUpdateRequest,
			encodeGRPCUpdateResponse,
			append(opts, grpc.ServerBefore(translateMetadataToContext()))...,
		),
	}
}

func NewGRPCClient(conn *grpcg.ClientConn, logger log.Logger) Service {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	opts := []grpc.ClientOption {}
	var authenticateEndpoint endpoint.Endpoint
	{
		authenticateEndpoint = grpc.NewClient(
			conn,
			"pb.Users",
			"Authenticate",
			grpcEncodeAuthenticateReqiest,
			grpcDecodeAuthenticateResponse,
			pbs.AuthenticateResponse{},
		).Endpoint()
		authenticateEndpoint = limiter(authenticateEndpoint)
	}
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = grpc.NewClient(
			conn,
			"pb.Users",
			"CreateUser",
			encodeGRPCCreateUserRequest,
			decodeGRPCCreateUserResponse,
			pbs.CreateUserResponse{},
		).Endpoint()
		createEndpoint = limiter(createEndpoint)
	}
	var findByEmailEndpoint endpoint.Endpoint
	{
		findByEmailEndpoint = grpc.NewClient(
			conn,
			"pb.Users",
			"FindUserByEmail",
			encodeGRPCFindByEmailRequest,
			decodeGRPCFindByEmailResponse,
			pbs.FindUserByEmailResponse{},
			append(opts, grpc.ClientBefore(translateJwtToMetadata()))...,
		).Endpoint()
		findByEmailEndpoint = limiter(findByEmailEndpoint)
	}
	var findByIdEndpoint endpoint.Endpoint
	{
		findByIdEndpoint = grpc.NewClient(
			conn,
			"pb.Users",
			"FindUserById",
			encodeGRPCFindByIdRequest,
			decodeGRPCFindByIdResponse,
			pbs.FindUserByIdResponse{},
			append(opts, grpc.ClientBefore(translateJwtToMetadata()))...,
		).Endpoint()
		findByIdEndpoint = limiter(findByIdEndpoint)
	}
	var updateEndpoint endpoint.Endpoint
	{
		updateEndpoint = grpc.NewClient(
			conn,
			"pb.Users",
			"UpdateUser",
			encodeGRPCUpdateRequest,
			decodeGRPCUpdateResponse,
			pbs.UpdateUserResponse{},
			append(opts, grpc.ClientBefore(translateJwtToMetadata()))...,
		).Endpoint()
		updateEndpoint = limiter(updateEndpoint)
	}

	return &Endpoints {
		AuthenticateEndpoint: authenticateEndpoint,
		CreateEndpoint: createEndpoint,
		FindByEmailEndpoint: findByEmailEndpoint,
		FindByIdEndpoint: findByIdEndpoint,
		UpdateEndpoint: updateEndpoint,
	}
}

func (t *grpcTransport) Authenticate(ctx context.Context, r *pbs.AuthenticateRequest) (*pbs.AuthenticateResponse, error) {
	_, resp, err := t.authenticate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.AuthenticateResponse), nil
}

func (t *grpcTransport) CreateUser(ctx context.Context, r *pbs.CreateUserRequest) (*pbs.CreateUserResponse, error) {
	_, resp, err := t.create.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.CreateUserResponse), nil
}

func (t *grpcTransport)	FindUserByEmail(ctx context.Context, r *pbs.FindUserByEmailRequest) (*pbs.FindUserByEmailResponse, error) {
	_, resp, err := t.findByEmail.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.FindUserByEmailResponse), nil
}

func (t *grpcTransport)	FindUserById(ctx context.Context, r *pbs.FindUserByIdRequest) (*pbs.FindUserByIdResponse, error) {
	_, resp, err := t.findById.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.FindUserByIdResponse), nil
}

func (t *grpcTransport)	UpdateUser(ctx context.Context, r *pbs.UpdateUserRequest) (*pbs.UpdateUserResponse, error) {
	_, resp, err := t.update.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pbs.UpdateUserResponse), nil
}

func grpcDecodeAuthenticateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.AuthenticateRequest)

	return &AuthenticateRequest {
		Login: req.Login,
		Password: req.Password,
	}, nil
}

func grpcEncodeAuthenticateReqiest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*AuthenticateRequest)

	return &pbs.AuthenticateRequest {
		Login: req.Login,
		Password: req.Password,
	}, nil
}

func grpcDecodeAuthenticateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.AuthenticateResponse)

	if resp.Err != "" {
		return &AuthenticateResponse {
			Err: resp.Err,
		}, nil
	}

	pbUser := protoUserToModelUser(*resp.User)
	return &AuthenticateResponse {
		User: &pbUser,
		Err: resp.Err,
	}, nil
}

func grpcEncodeAuthenticateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*AuthenticateResponse)

	if resp.Err != "" {
		return &pbs.AuthenticateResponse {
			Err: resp.Err,
		}, nil
	}

	user := modelsUserToProtoUser(*resp.User)
	return &pbs.AuthenticateResponse {
		User: &user,
	}, nil
}

func decodeGRPCCreateUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.CreateUserRequest)
	
	user := protoUserToModelUser(*req.User)
	return &CreateUserRequest {
		User: &user,
	}, nil
}

func encodeGRPCCreateUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*CreateUserResponse)

	return &pbs.CreateUserResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func decodeGRPCFindByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.FindUserByEmailRequest)

	return &FindUserByEmailRequest {
		Email: req.Email,
	}, nil
}

func encodeGRPCFindByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*FindUserByEmailResponse)

	if resp.Err != "" {
		return &pbs.FindUserByEmailResponse {
			Err: resp.Err,
		}, nil
	}

	pbUser := modelsUserToProtoUser(*resp.User)
	return &pbs.FindUserByEmailResponse {
		User: &pbUser,
		Err: resp.Err,
	}, nil
}

func decodeGRPCFindByIdRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.FindUserByIdRequest)

	return &FindUserByIdRequest {
		Id: int(req.Id),
	}, nil
}

func encodeGRPCFindByIdResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*FindUserByIdResponse)

	if resp.Err != "" {
		return &pbs.FindUserByIdResponse {
			Err: resp.Err,
		}, nil
	}

	pbUser := modelsUserToProtoUser(*resp.User)
	return &pbs.FindUserByIdResponse {
		User: &pbUser,
		Err: resp.Err,
	}, nil
}

func decodeGRPCUpdateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pbs.UpdateUserRequest)

	user := protoUserToModelUser(*req.User)
	return &UpdateUserRequest {
		User: &user,
	}, nil
}

func encodeGRPCUpdateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*UpdateUserResponse)
	if resp.Err != "" {
		return &pbs.UpdateUserResponse {
			Err: resp.Err,
		}, nil
	}

	return &pbs.UpdateUserResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func encodeGRPCCreateUserRequest(_context context.Context, r interface{}) (interface{}, error) {
	req := r.(*CreateUserRequest)
	
	pbUser := modelsUserToProtoUser(*req.User)
	return &pbs.CreateUserRequest {
		User: &pbUser,
	}, nil
}
			
func decodeGRPCCreateUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.CreateUserResponse)
	
	return &CreateUserResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func encodeGRPCFindByEmailRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*FindUserByEmailRequest)

	return &pbs.FindUserByEmailRequest {
		Email: req.Email,
	}, nil
}

func decodeGRPCFindByEmailResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.FindUserByEmailResponse)

	if resp.Err != "" {
		return &FindUserByEmailResponse {
			Err: resp.Err,
		}, nil
	}

	user := protoUserToModelUser(*resp.User)
	return &FindUserByEmailResponse {
		User: &user,
		Err: resp.Err,
	}, nil
}

func encodeGRPCFindByIdRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*FindUserByIdRequest)

	return &pbs.FindUserByIdRequest {
		Id: int32(req.Id),
	}, nil
}
			
func decodeGRPCFindByIdResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.FindUserByIdResponse)

	if resp.Err != "" {
		return &FindUserByIdResponse {
			Err: resp.Err,
		}, nil
	}

	user := protoUserToModelUser(*resp.User)
	return &FindUserByIdResponse {
		User: &user,
	}, nil
}

func encodeGRPCUpdateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*UpdateUserRequest)

	pbUser := modelsUserToProtoUser(*req.User)
	return &pbs.UpdateUserRequest {
		User: &pbUser,
	}, nil
}
func decodeGRPCUpdateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pbs.UpdateUserResponse)
	if resp.Err != "" {
		return &UpdateUserResponse {
			Err: resp.Err,
		}, nil
	}

	return &UpdateUserResponse {
		Ok: resp.Ok,
		Err: resp.Err,
	}, nil
}

func modelsUserToProtoUser(p User) pbs.User {
	return pbs.User {
		Id: int32(p.ID),
		Email: p.Email,
		Login: p.Login,
		Contacts: p.Contacts,
		CreatedAt: p.CreatedAt,
		Role: p.Role,
		Token: p.Token,
		Password: p.Password,
	}
}

func protoUserToModelUser(pb pbs.User) User {
	return User{
		ID: int(pb.Id),
		Email: pb.Email,
		Login: pb.Login,
		Contacts: pb.Contacts,
		CreatedAt: pb.CreatedAt,
		Token: pb.Token,
		Role: pb.Role,
		Password: pb.Password,
	}
}