package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	errLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
	cocain "github.com/inhumanLightBackend/cocain/logic"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	config := createConfig()
	var (
		httpAddr     = config.HttPort
		consulAddr   = config.ConsulPort
		retryMax     = config.MaxRetry
		retryTimeout = time.Duration(config.MaxTimeout) * time.Millisecond
		db           = createDatabase(config.DatabaseUrl)
	)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if len(consulAddr) > 0 {
			consulConfig.Address = consulAddr
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	// Transport domain.
	// _ := stdopentracing.GlobalTracer() // no-op
	// _, _ := stdzipkin.NewTracer(nil, stdzipkin.WithNoopTracer(true))
	// ctx := context.Background()

	repo := cocain.NewGateRepo(db, logger)
	gates := cocain.NewGateService(repo, logger)
	middleware := NewMiddleware(gates)

	r := mux.NewRouter()
	handlers := NewHandlerFactory(client, logger, retryMax, retryTimeout)
	r.HandleFunc("/signin", handlers.SignInEndpoint(gates)).Methods("POST")
	r.HandleFunc("/logout", handlers.LogoutEndpoint(gates)).Methods("GET")
	
	prefix := "/v1"
	auth := r.PathPrefix(prefix).Subrouter()
	auth.Use(middleware.Translate)
	auth.PathPrefix("/auth").Handler(http.StripPrefix(fmt.Sprintf("%s/auth", prefix), handlers.AuthHandler()))
	auth.PathPrefix("/user").Handler(http.StripPrefix(fmt.Sprintf("%s/user", prefix), handlers.UserHandler()))
	auth.PathPrefix("/notifications").Handler(http.StripPrefix(fmt.Sprintf("%s/notifications", prefix), handlers.NotificationsHandler()))
	auth.PathPrefix("/support").Handler(http.StripPrefix(fmt.Sprintf("%s/support", prefix), handlers.SupportHandler()))

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "http", "addr", httpAddr)
		errs <- http.ListenAndServe(httpAddr, r)
	}()

	logger.Log("Terminate", <-errs)
}

func createConfig() *cocain.Config {
	var configPath string
	flag.StringVar(&configPath, "config-path", "cocain/logic/_config.toml", "path to config file")
	flag.Parse()
	config := cocain.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		errLog.Fatal(err)
	}

	return config
}

func createDatabase(databaseURL string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		errLog.Fatal(err)
	}

	ddl, err := ioutil.ReadFile("cocain/schema/init_schema.sql")
	if err != nil {
		errLog.Fatal(err)
	}

	_, err = db.Exec(string(ddl))
	if err != nil {
		errLog.Fatal(err)
	}

	return db
}
