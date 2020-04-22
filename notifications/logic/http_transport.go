package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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

	createHandler := kithttp.NewServer(
		ep.CreateEndpoint,
		httpDecodeCreateRequest,
		httpEncodeCreateResponse,
		opts...,
	)

	findByIdHandler := kithttp.NewServer(
		ep.FindByIdEndpoint,
		httpDecodeFindByIdRequest,
		httpEncodeFindByIdResponse,
		opts...,
	)

	checkHandler := kithttp.NewServer(
		ep.CheckEndpoint,
		httpDecodeCheckRequest,
		httpEncodeCheckResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/create", createHandler).Methods("POST")
	r.Handle("/findById", findByIdHandler).Methods("GET")
	r.Handle("/check", checkHandler).Methods("POST")

	return accessControl(r)
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

func httpDecodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := &Notification{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, errBadRequest
	}
	return &CreateRequest {
		Notification: req,
	}, nil
}

func httpEncodeCreateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*CreateResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeFindByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &FindByIdRequest{}, nil
}

func httpEncodeFindByIdResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*FindByIdResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := &CheckRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, errBadRequest
	}
	return req, nil
}

func httpEncodeCheckResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*CheckResponse)
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