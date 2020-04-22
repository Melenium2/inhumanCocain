package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AuthenticateEndpoint endpoint.Endpoint
	CreateEndpoint       endpoint.Endpoint
	FindByEmailEndpoint  endpoint.Endpoint
	FindByIdEndpoint     endpoint.Endpoint
	UpdateEndpoint       endpoint.Endpoint
}

type (
	AuthenticateRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	AuthenticateResponse struct {
		User *User  `json:"user,omitempty"`
		Err  string `json:"error,omitempty"`
	}

	CreateUserRequest struct {
		User *User `json:"user"`
	}

	CreateUserResponse struct {
		Ok  string `json:"ok,omitempty"`
		Err string `json:"error,omitempty"`
	}

	FindUserByEmailRequest struct {
		Email string `json:"email"`
	}

	FindUserByEmailResponse struct {
		User *User  `json:"user,omitempty"`
		Err  string `json:"error,omitempty"`
	}

	FindUserByIdRequest struct {
		Id int `json:"id"`
	}

	FindUserByIdResponse struct {
		User *User  `json:"user,omitempty"`
		Err  string `json:"error,omitempty"`
	}

	UpdateUserRequest struct {
		User *User `json:"user"`
	}

	UpdateUserResponse struct {
		Ok  string `json:"message,omitempty"`
		Err string `json:"error,omitempty"`
	}
)

func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		AuthenticateEndpoint: MakeAuthnticateEndpoint(s),
		CreateEndpoint:       MakeCreateEndpoint(s),
		FindByEmailEndpoint:  MakeFindByEmailEndpoint(s),
		FindByIdEndpoint:     MakeFindByIdEndpoint(s),
		UpdateEndpoint:       MakeUpdateEndpoint(s),
	}
}

func (ep *Endpoints) Authenticate(ctx context.Context, login string, pass string) (*User, error) {
	r, err := ep.AuthenticateEndpoint(ctx, &AuthenticateRequest{
		Login: login,
		Password: pass,
	})
	if err != nil {
		return nil, err
	}

	resp := r.(*AuthenticateResponse)
	return resp.User, nil
}

func (ep *Endpoints) CreateUser(ctx context.Context, u *User) (string, error) {
	r, err := ep.CreateEndpoint(ctx, &CreateUserRequest{
		User: u,
	})
	if err != nil {
		return "", err
	}
	response := r.(*CreateUserResponse)
	return response.Ok, nil
}

func (ep *Endpoints) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	r, err := ep.FindByEmailEndpoint(ctx, &FindUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	response := r.(*FindUserByEmailResponse)
	return response.User, nil
}

func (ep *Endpoints) FindUserById(ctx context.Context, id int) (*User, error) {
	r, err := ep.FindByIdEndpoint(ctx, &FindUserByIdRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	response := r.(*FindUserByIdResponse)
	return response.User, nil
}

func (ep *Endpoints) UpdateUser(ctx context.Context, u *User) (string, error) {
	r, err := ep.UpdateEndpoint(ctx, &UpdateUserRequest{
		User: u,
	})
	if err != nil {
		return "", err
	}
	reponse := r.(*UpdateUserResponse)
	println(reponse.Ok, reponse.Err)
	return reponse.Ok, nil
}

func MakeAuthnticateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*AuthenticateRequest)
		u, err := s.Authenticate(ctx, req.Login, req.Password)
		if err != nil {
			println(err.Error())
			return &AuthenticateResponse{
				Err: err.Error(),
			}, nil
		}

		return &AuthenticateResponse{
			User: u,
		}, nil
	}
}

func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CreateUserRequest)
		mess, err := s.CreateUser(ctx, req.User)
		if err != nil {
			println(err.Error())
			return &CreateUserResponse{
				Err: err.Error(),
			}, nil
		}

		return &CreateUserResponse{
			Ok: mess,
		}, nil
	}
}

func MakeFindByEmailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*FindUserByEmailRequest)
		user, err := s.FindUserByEmail(ctx, req.Email)
		if err != nil {
			return &FindUserByEmailResponse{
				Err: err.Error(),
			}, nil
		}

		return &FindUserByEmailResponse{
			User: user,
		}, nil
	}
}

func MakeFindByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*FindUserByIdRequest)
		user, err := s.FindUserById(ctx, req.Id)
		if err != nil {
			return &FindUserByIdResponse{
				Err: err.Error(),
			}, nil
		}

		return &FindUserByIdResponse{
			User: user,
		}, nil
	}
}

func MakeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UpdateUserRequest)
		mess, err := s.UpdateUser(ctx, req.User)
		if err != nil {
			println(err.Error())
			return &UpdateUserResponse{
				Err: err.Error(),
			}, nil
		}
		return &UpdateUserResponse{
			Ok: mess,
		}, nil
	}
}
