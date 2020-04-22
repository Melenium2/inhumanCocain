package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	auth "github.com/inhumanLightBackend/auth/logic"
	"github.com/inhumanLightBackend/cocain"
	gates "github.com/inhumanLightBackend/cocain/logic"
	notifications "github.com/inhumanLightBackend/notifications/logic"
	support "github.com/inhumanLightBackend/support/logic"
	user "github.com/inhumanLightBackend/user/logic"
)

var (
	errBadRequest = errors.New("Bad request")
)

type HandlerFactory struct {
	consul       consul.Client
	logger       log.Logger
	maxRetry     int
	retryTimeout time.Duration
}

func NewHandlerFactory(client consul.Client, logger log.Logger, maxRetry int, retryTimeout time.Duration) *HandlerFactory {
	return &HandlerFactory{
		consul:       client,
		logger:       logger,
		maxRetry:     maxRetry,
		retryTimeout: retryTimeout,
	}
}

func (hf *HandlerFactory) AuthHandler() http.Handler {
	var (
		tags        = []string{}
		passingOnly = true
		endpoints   = &auth.Endpoints{}
		instancer   = consulsd.NewInstancer(hf.consul, hf.logger, "auth", tags, passingOnly)
	)
	{
		factory := authFactory(auth.MakeGenerateEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.GenerateEndpoint = retry
	}
	{
		factory := authFactory(auth.MakeValidateEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.ValidateEndpoint = retry
	}

	return auth.NewHTTPTransport(endpoints, hf.logger)
}

func (hf *HandlerFactory) UserHandler() http.Handler {
	var (
		tags        = []string{}
		passingOnly = true
		endpoints   = &user.Endpoints{}
		instancer   = consulsd.NewInstancer(hf.consul, hf.logger, "user", tags, passingOnly)
	)
	{
		factory := userFactory(user.MakeAuthnticateEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.AuthenticateEndpoint = retry
	}
	{
		factory := userFactory(user.MakeCreateEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.CreateEndpoint = retry
	}
	{
		factory := userFactory(user.MakeFindByEmailEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.FindByEmailEndpoint = retry
	}
	{
		factory := userFactory(user.MakeFindByIdEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.FindByIdEndpoint = retry
	}
	{
		factory := userFactory(user.MakeUpdateEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.UpdateEndpoint = retry
	}

	return user.NewHTTPTransport(endpoints, hf.logger)
}

func (hf *HandlerFactory) NotificationsHandler() http.Handler {
	var (
		tags        = []string{}
		passingOnly = true
		endpoints   = &notifications.Endpoints{}
		instancer   = consulsd.NewInstancer(hf.consul, hf.logger, "notifications", tags, passingOnly)
	)
	{
		factory := notifFactoty(notifications.MakeCreateEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.CreateEndpoint = retry
	}
	{
		factory := notifFactoty(notifications.MakeFindByIdEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.FindByIdEndpoint = retry
	}
	{
		factory := notifFactoty(notifications.MakeCheckEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.CheckEndpoint = retry
	}

	return notifications.NewHTTPTransport(endpoints, hf.logger)
}

func (hf *HandlerFactory) SupportHandler() http.Handler {
	var (
		tags        = []string{}
		passingOnly = true
		endpoints   = &support.Endpoints{}
		instancer   = consulsd.NewInstancer(hf.consul, hf.logger, "support", tags, passingOnly)
	)
	{
		factory := supportFactory(support.MakeCreateTicketEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.CreateTicketEndpoint = retry
	}
	{
		factory := supportFactory(support.MakeGetTicketEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.GetTicketEndpoint = retry
	}
	{
		factory := supportFactory(support.MakeTicketsEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.TicketsEndpoint = retry
	}
	{
		factory := supportFactory(support.MakeAcceptTicketEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.AcceptTicketEndpoint = retry
	}
	{
		factory := supportFactory(support.MakeAddMessageEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.AddMessageEndpoint = retry
	}
	{
		factory := supportFactory(support.MakeGetMessagesEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.GetMessagesEndpoint = retry
	}
	{
		factory := supportFactory(support.MakeChangeStatusEndpoint, hf.logger)
		endpointer := sd.NewEndpointer(instancer, factory, hf.logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(hf.maxRetry, hf.retryTimeout, balancer)
		endpoints.ChangeStatusEndpoint = retry
	}

	return support.NewHTTPTransport(endpoints, hf.logger)
}

// Edpoint (/signin) for user login
func (hf *HandlerFactory) SignInEndpoint(s gates.Gates) http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err:= json.NewDecoder(r.Body).Decode(req); err != nil {
			SendError(w, r, http.StatusBadRequest, errBadRequest)
		}

		res, err := s.SignIn(r.Context(), req.Login, req.Password)
		if err != nil {
			SendError(w, r, http.StatusBadRequest, err)
		}

		Respond(w, r, http.StatusOK, res)
	}
}

// Endpoint (/logout) for user logout
func (hf *HandlerFactory) LogoutEndpoint(s gates.Gates) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("X-Opaque")
		if header == "" {
			SendError(w, r, http.StatusUnauthorized, errNotAuthorized)
			return
		}
		if err := s.Logout(r.Context(), header); err != nil {
			SendError(w, r, http.StatusInternalServerError, err)
			return 
		}

		Respond(w, r, http.StatusOK, map[string]string {
			"message": "ok",
		})
	}
}

func authFactory(makeEndpoint func(auth.Service) endpoint.Endpoint, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, conn, err := cocain.AuthGRPCService(instance, logger)
		if err != nil {
			return nil, nil, err
		}
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}

func userFactory(makeEndpoint func(user.Service) endpoint.Endpoint, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		s, c, err := cocain.UserGRPCService(instance, logger)
		if err != nil {
			return nil, nil, err
		}
		endpoint := makeEndpoint(s)

		return endpoint, c, nil
	}
}

func notifFactoty(makeEndpoint func(notifications.Service) endpoint.Endpoint, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		s, c, err := cocain.NotificationsGRPCService(instance, logger)
		if err != nil {
			return nil, nil, err
		}
		endpoint := makeEndpoint(s)

		return endpoint, c, nil
	}
}

func supportFactory(makeEndpoint func(support.Service) endpoint.Endpoint, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		s, c, err := cocain.SupportGRPCService(instance, logger)
		if err != nil {
			return nil, nil, err
		}
		endpoint := makeEndpoint(s)

		return endpoint, c, nil
	}
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{})  {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func SendError(w http.ResponseWriter, r *http.Request, code int, err error) {
	Respond(w, r, code, map[string]string{"error": err.Error()})
}


