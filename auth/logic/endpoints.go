package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GenerateEndpoint endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

type GenerateRequest struct {
	UserId int32  `json:"userId"`
	Role   string `json:"role"`
}

type GenerateResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"error,omitempty"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	UserId int32  `json:"userId,omitempty"`
	Role   string `json:"role,omitempty"`
	Token  string `json:"token,omitempty"`
	Err    string `json:"error,omitempty"`
}

func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		GenerateEndpoint: MakeGenerateEndpoint(s),
		ValidateEndpoint: MakeValidateEndpoint(s),
	}
}

func (ep *Endpoints) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	r, err := ep.GenerateEndpoint(ctx, req)
	if err != nil {
		println(err)
		return nil, err
	}

	response := r.(*GenerateResponse)
	return response, nil
}

func (ep *Endpoints) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	r, err := ep.ValidateEndpoint(ctx, req)
	if err != nil {
		println(err)
		return nil, err
	}

	response := r.(*ValidateResponse)
	return response, nil
}

func MakeGenerateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*GenerateRequest)
		mess, err := s.Generate(ctx, req)
		if err != nil {
			println(err.Error())
			return &GenerateResponse{
				Err: err.Error(),
			}, nil
		}

		return &GenerateResponse{
			Token: mess.Token,
			Err:   mess.Err,
		}, nil
	}
}

func MakeValidateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*ValidateRequest)
		mess, err := s.Validate(ctx, req)
		if err != nil {
			return &ValidateResponse {
				Err: err.Error(),
			}, nil
		}

		return mess, nil
	}
}