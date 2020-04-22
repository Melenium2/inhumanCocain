package support

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	errBadRoute        = errors.New("Bad route")
	errBadRequest      = errors.New("Bad request")
	errNotFound        = errors.New("Not found")
	errPermissionsDeni = errors.New("Permissions denied")
)

func NewHTTPTransport(ep *Endpoints, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(log.With(logger, "tag", "http")),
		kithttp.ServerErrorEncoder(encodeHttpError),
	}

	createTicketHandler := kithttp.NewServer(
		ep.CreateTicketEndpoint,
		httpDecodeCreateTicketRequest,
		httpEncodeCreateTicketResponse,
		opts...,
	)
	getTicketHandler := kithttp.NewServer(
		ep.GetTicketEndpoint,
		httpDecodeGetTicketRequest,
		httpEncodeGetTicketResponse,
		opts...,
	)
	ticketsHandler := kithttp.NewServer(
		ep.TicketsEndpoint,
		httpDecodeTicketsRequest,
		httpEncodeTicketsResponse,
		opts...,
	)
	acceptTicketHandler := kithttp.NewServer(
		ep.AcceptTicketEndpoint,
		httpDecodeAcceptTicketRequest,
		httpEncodeAcceptTicketResponse,
		opts...,
	)
	addMessageHandler := kithttp.NewServer(
		ep.AddMessageEndpoint,
		httpDecodeAddMessageRequest,
		httpEncodeAddMessageResponse,
		opts...,
	)
	getMessagesHandler := kithttp.NewServer(
		ep.GetMessagesEndpoint,
		httpDecodeGetMessagesRequest,
		httpEncodeGetMessagesResponse,
		opts...,
	)
	changeStatusHandler := kithttp.NewServer(
		ep.ChangeStatusEndpoint,
		httpDecodeChangeStatusRequest,
		httpEncodeChangeStatusResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/ticket/create", createTicketHandler).Methods("POST")
	r.Handle("/ticket/get", getTicketHandler).Methods("GET")
	r.Handle("/ticket/all", ticketsHandler).Methods("GET")
	r.Handle("/ticket/accept", acceptTicketHandler).Methods("GET")
	r.Handle("/message/add", addMessageHandler).Methods("POST")
	r.Handle("/message/get", getMessagesHandler).Methods("GET")
	r.Handle("/ticket/status", changeStatusHandler).Methods("GET")

	return accessControl(r)
}

func httpDecodeCreateTicketRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := &Ticket{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, errBadRequest
	}
	println("httpDecodeCreateTicketRequest")
	return &CreateTicketRequest{
		Ticket: req,
	}, nil
}

func httpEncodeCreateTicketResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*CreateTicketResponse)
	println("httpEncodeCreateTicketResponse", resp.Err, resp.Ok)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeGetTicketRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		ID  int
		err error
	)
	{
		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			return nil, errBadRoute
		}
		ID, err = strconv.Atoi(id[0])
		if err != nil {
			return nil, errBadRoute
		}
	}

	return &GetTicketRequest{
		TicketId: ID,
	}, nil
}

func httpEncodeGetTicketResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*GetTicketResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeTicketsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return &TicketsRequest {}, nil
}

func httpEncodeTicketsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*TicketsResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeAcceptTicketRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		ID  int
		err error
	)
	{
		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			return nil, errBadRoute
		}
		ID, err = strconv.Atoi(id[0])
		if err != nil {
			return nil, errBadRoute
		}
	}
	
	return &AcceptTicketRequest {
		TicketId: ID,
	}, nil
}

func httpEncodeAcceptTicketResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*AcceptTicketResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeAddMessageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := &TicketMessage{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, errBadRequest
	}
	
	return &AddMessageRequest {
		Message: req,
	}, nil
}

func httpEncodeAddMessageResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*AddMessageResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeGetMessagesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		ID  int
		err error
	)
	{
		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			return nil, errBadRoute
		}
		ID, err = strconv.Atoi(id[0])
		if err != nil {
			return nil, errBadRoute
		}
	}
	
	return &GetMessagesRequest {
		TicketId: ID,
	}, nil
}

func httpEncodeGetMessagesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*GetMessagesResponse)
	if resp.Err == "" {
		return encodeHTTPResponse(ctx, http.StatusOK, w, response)
	}
	encodeHttpError(ctx, getHTTPError(resp.Err), w)
	return nil
}

func httpDecodeChangeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var (
		ID  int
		status string
		err error
	)
	{
		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			return nil, errBadRoute
		}
		ID, err = strconv.Atoi(id[0])
		if err != nil {
			return nil, errBadRoute
		}
		status = r.URL.Query().Get("status")
		if status == "" {
			return nil, errBadRoute
		}
	}
	
	return &ChangeStatusRequest {
		TicketId: ID,
		Status: status,
	}, nil
}

func httpEncodeChangeStatusResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*ChangeStatusResponse)
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
	e := json.NewEncoder(w).Encode(map[string]interface{}{
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
