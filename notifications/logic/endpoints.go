package notifications

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateEndpoint   endpoint.Endpoint
	FindByIdEndpoint endpoint.Endpoint
	CheckEndpoint    endpoint.Endpoint
}

type (
	CreateRequest struct {
		Notification *Notification `json:"notification"`
	}

	CreateResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}

	FindByIdRequest struct {
	}

	FindByIdResponse struct {
		Notifications []*Notification `json:"notifications,omitempty"`
		Err           string          `json:"error,omitempty"`
	}

	CheckRequest struct {
		Indexes []int `json:"indexes"`
	}

	CheckResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}
)

func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		CreateEndpoint: MakeCreateEndpoint(s),
		FindByIdEndpoint: MakeFindByIdEndpoint(s),
		CheckEndpoint: MakeCheckEndpoint(s),
	}
}

func (ep *Endpoints) Create(ctx context.Context, n *Notification) (string, error) {
	r, err := ep.CreateEndpoint(ctx, &CreateRequest{
		Notification: n,	
	})
	if err != nil {
		return "", err
	}

	resp := r.(*CreateResponse)
	return resp.Ok, errors.New(resp.Err)
}

func (ep *Endpoints) FindById(ctx context.Context) ([]*Notification, error) {
	r, err := ep.FindByIdEndpoint(ctx, &FindByIdRequest{})
	if err != nil {
		return nil, err
	}

	resp := r.(*FindByIdResponse)
	return resp.Notifications, errors.New(resp.Err)
}

func (ep *Endpoints) Check(ctx context.Context, indexes []int) (string, error) {
	r, err := ep.CheckEndpoint(ctx, &CheckRequest{
		Indexes: indexes,
	})
	if err != nil {
		return "", err
	}

	resp := r.(*CheckResponse)
	return resp.Ok, errors.New(resp.Err)
}

func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CreateRequest)
		n, err := s.Create(ctx, req.Notification)

		if err != nil {
			return &CreateResponse {
				Err: err.Error(),
			}, nil
		}
		
		return &CreateResponse{
			Ok: n,
		}, nil
	}
}	

func MakeFindByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ns, err := s.FindById(ctx)
		
		if err != nil {
			return &FindByIdResponse {
				Err: err.Error(),
			}, nil
		}

		return &FindByIdResponse{
			Notifications: ns,
		}, nil
	}
}

func MakeCheckEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CheckRequest)
		ok, err := s.Check(ctx, req.Indexes)

		if err != nil {
			return &CheckResponse {
				Err: err.Error(),
			}, nil
		}

		return &CheckResponse{
			Ok: ok,
		}, nil
	}
}



