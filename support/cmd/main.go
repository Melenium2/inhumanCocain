package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	kitlog "github.com/go-kit/kit/log"
	support "github.com/inhumanLightBackend/support/logic"
	pbs "github.com/inhumanLightBackend/support/pb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	logger := createLogger()
	l := kitlog.With(logger, "tag", "start")
	l.Log("msg", "created logger")

	config, err := createConfig()
	if err != nil {
		l.Log("msg", "can not create config", "err", err)
		os.Exit(-1)
	}

	db, err := createDatabase(config.DatabaseUrl)
	if err != nil {
		l.Log("msg", "failed to connect to database", "err", err)
		os.Exit(-1)
	}
	defer db.Close()
	l.Log("msg", "connect to database")
	repo := support.NewRepository(db, logger)
	service := support.NewService(repo, logger)
	ep := support.NewEndpoints(service)

	errs := make(chan error)

	go func() {
		grpcTransport := support.NewGRPCServer(ep, logger)
		l.Log("msg", "created grpc transport")

		listner, err := net.Listen("tcp", config.GrpcPort)
		if err != nil {
			errs <- err
			return
		}
		l.Log("started", "grpc", "listner", config.GrpcPort, "msg", "listening")
		grpcServer := grpc.NewServer()
		pbs.RegisterSupportServer(grpcServer, grpcTransport)
		l.Log("transport", "grpc", "address", config.GrpcPort, "msg", "listening")
		errs <- grpcServer.Serve(listner)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <-errs)
}

func createConfig() (*support.Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config-path", "support/logic/_config.toml", "path to config file")
	flag.Parse()
	config := support.NewConfig()
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

	ddl, err := ioutil.ReadFile("support/schema/init_schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(ddl))

	return db, err
}