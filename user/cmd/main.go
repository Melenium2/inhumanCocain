package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	kitlog "github.com/go-kit/kit/log"
	user "github.com/inhumanLightBackend/user/logic"
	pbs "github.com/inhumanLightBackend/user/pb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	logger := createLogger()
	startLogger := kitlog.With(logger, "tag", "start")
	startLogger.Log("msg", "created logger")
	
	config, err := createConfig()
	if err != nil {
		startLogger.Log("msg", "can not create config", "err", err)
		os.Exit(-1)
	}

	db, err := createDatabase(config.DatabaseURL)
	if err != nil {
		startLogger.Log("msg", "failed to connect to database", "err", err)
		os.Exit(-1)
	}
	defer db.Close()
	startLogger.Log("msg", "connect to database")
	repo := user.NewRepo(db)
	service := user.NewService(repo) 
	ep := user.NewEndpoints(service)

	errs := make(chan error)

	go func() {
		httpTransport := user.NewHTTPTransport(ep, logger)
		startLogger.Log("msg", "created http transport")
		startLogger.Log("transport", "http", "address", config.HTTPPort, "msg", "listening")
		errs <- http.ListenAndServe(config.HTTPPort, httpTransport)
	}()

	go func() {
		grpcTransport := user.NewGRPCServer(ep, logger)
		startLogger.Log("msg", "created grpc transport")
		
		listner, err := net.Listen("tcp", config.GRPCPort)
		if err != nil {
			errs <- err
			return
		}
		startLogger.Log("started", "grpc", "listner", config.GRPCPort, "msg", "listening")
		grpcServer := grpc.NewServer()
		pbs.RegisterUsersServer(grpcServer, grpcTransport)
		startLogger.Log("transport", "grpc", "address", config.GRPCPort, "msg", "listening")
		errs <- grpcServer.Serve(listner)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <-errs)
}

func createConfig() (*user.Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config-path", "user/logic/_config.toml", "path to config file")
	flag.Parse()
	config := user.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func createLogger() kitlog.Logger {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	return kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC())
}

func createDatabase(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	ddl, err := ioutil.ReadFile("user/schema/create_schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(ddl))

	return db, err
}