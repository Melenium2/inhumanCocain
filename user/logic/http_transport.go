package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	errBadRoute = errors.New("Bad route")
	errBadRequest = errors.New("Bad request")
	errNotFound = errors.New("Not found")
	errPermissionsDeni = errors.New("Permissions denied")
)

func NewHTTPTransport(ep *Endpoints, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption {
		kithttp.ServerErrorLogger(log.With(logger, "tag", "http")),
		kithttp.ServerErrorEncoder(encodeHttpError),
	}

	authenticateHandler := kithttp.NewServer(
		ep.AuthenticateEndpoint,
		httpDecodeAuthenticateRequest,
		httpEncodeAuthenticateResponse,
		opts...,
	)

	createHandler := kithttp.NewServer(
		ep.CreateEndpoint,
		decodeHTTPCreateUserRequest,
		encodeHTTPCreateUserResponse,
		opts...,
	)

	findByEmailHandler := kithttp.NewServer(
		ep.FindByEmailEndpoint,
		decodeHTTPFindByUserEmailRequest,
		encodeHTTPFindByUserEmailResponse,
		opts...,
	)

	findByIdHandler := kithttp.NewServer(
		ep.FindByIdEndpoint,
		decodeHTTPFindByUserIdRequest,
		encodeHTTPFindByUserIdResponse,
		opts...,
	)

	updateHandler := kithttp.NewServer(
		ep.UpdateEndpoint,
		decodeHTTPUpdateUserRequest,
		encodeHTTPUpdateUserResponse,
		opts...,
	)

	r := mux.NewRouter()
	main := r.PathPrefix("").Subrouter()
	main.Handle("/auth", authenticateHandler).Methods("POST")
	main.Handle("/create", createHandler).Methods("POST")
	main.Handle("/findByEmail", findByEmailHandler).Methods("GET")
	main.Handle("/findById", findByIdHandler).Methods("GET")
	main.Handle("/update", updateHandler).Methods("POST")

	return accessControl(main)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func httpDecodeAuthenticateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := &AuthenticateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, errBadRequest
	}
	return req, nil
}

func httpEncodeAuthenticateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*AuthenticateResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func decodeHTTPCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return nil, errBadRequest
	}
	println(fmt.Sprintf("%v", user))
	return &CreateUserRequest {
		User: user,
	}, nil
}

func encodeHTTPCreateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*CreateUserResponse)
	if resp.Err == "" && resp.Ok != "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func decodeHTTPFindByUserEmailRequest(_ context.Context, r *http.Request) (interface{}, error) {
	email, ok := r.URL.Query()["email"]
	if !ok && len(email[0]) == 0 {
		return nil, errBadRoute
	}

	return &FindUserByEmailRequest {
		Email: email[0],
	}, nil
}

func encodeHTTPFindByUserEmailResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*FindUserByEmailResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func decodeHTTPFindByUserIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, ok := r.URL.Query()["id"]
	if !ok && len(id[0]) == 0 {
		return nil, errBadRoute
	}

	ID, err := strconv.Atoi(id[0])
	if err != nil {
		return nil, errBadRoute
	}

	return &FindUserByIdRequest {
		Id: ID,
	}, nil
}

func encodeHTTPFindByUserIdResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*FindUserByIdResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func decodeHTTPUpdateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return nil, errBadRequest
	}

	return &UpdateUserRequest {
		User: user,
	}, nil
}

func encodeHTTPUpdateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*UpdateUserResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func encodeHTTPResponse(_ context.Context, code int, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(code)
	if response == nil {
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeHttpError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case errBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	case errBadRoute:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	e := json.NewEncoder(w).Encode(map[string]interface{} {
		"error": err.Error(),
	})
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getHTTPError(str string) error {
	if str == "not found" {
		return errNotFound
	}

	return errors.New(str)
}