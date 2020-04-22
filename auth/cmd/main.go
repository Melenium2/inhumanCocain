package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	kitlog "github.com/go-kit/kit/log"
	auth "github.com/inhumanLightBackend/auth/logic"
	pbs "github.com/inhumanLightBackend/auth/pb"
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

	startLogger.Log("msg", "connect to database")
	service := auth.NewService(config) 
	ep := auth.NewEndpoints(service)

	errs := make(chan error)

	go func() {
		grpcTransport := auth.NewGRPCServer(ep, logger)
		startLogger.Log("msg", "created grpc transport")
		
		listner, err := net.Listen("tcp", config.GrpcPort)
		if err != nil {
			errs <- err
			return
		}
		startLogger.Log("started", "grpc", "listner", config.GrpcPort, "msg", "listening")
		grpcServer := grpc.NewServer()
		pbs.RegisterAuthServer(grpcServer, grpcTransport)
		startLogger.Log("transport", "grpc", "address", config.GrpcPort, "msg", "listening")
		errs <- grpcServer.Serve(listner)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <-errs)
}

func createConfig() (*auth.Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config-path", "auth/logic/_config.toml", "path to config file")
	flag.Parse()
	config := auth.NewConfig()
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
