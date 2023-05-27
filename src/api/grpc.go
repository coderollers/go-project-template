package api

import (
	"context"
	"fmt"
	"net"

	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	grpcServices "my-microservice/api/grpc"
	"my-microservice/configuration"
	"my-microservice/protos"
)

func StartGrpc(ctx context.Context) (*grpc.Server, *grpcweb.WrappedGrpcServer) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	// Set up grpc
	log.Debugf("Setting up GRPC")
	grpcServer := grpc.NewServer()

	// Example GRPC service
	protos.RegisterGreeterServer(grpcServer, &grpcServices.GreeterService{})
	// TEMPLATE: Register GRPC services

	if conf.HttpPort == conf.GrpcPort {
		// Ports match, we return a GRPC-Web wrapper to use with our regular HTTP listener
		log.Infof("Multiplexed GRPC native and GRPC-Web mode on :%d", conf.GrpcPort)
		return grpcServer, grpcweb.WrapServer(grpcServer)
	}

	// Start native GRPC endpoint if ports differ
	go nativeGrpc(ctx, grpcServer)

	return nil, nil
}

func nativeGrpc(ctx context.Context, grpcServer *grpc.Server) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	// Native GRPC endpoint
	log.Infof("GRPC native mode on :%d", conf.GrpcPort)

	var (
		listener net.Listener
		err      error
	)

	// Set up the listener
	if listener, err = net.Listen("tcp", fmt.Sprintf(":%d", conf.GrpcPort)); err != nil {
		log.Fatalf("Cannot open listener socket: %s", err.Error())
	}

	// Start the HTTP Server
	go func() {
		concurrency.GlobalWaitGroup.Add(1)
		defer concurrency.GlobalWaitGroup.Done()

		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve GRPC endpoint: %s", err.Error())
		}
	}()

	// Block until SIGTERM/SIGINT
	<-ctx.Done()

	// Clean up and shutdown the HTTP server
	grpcServer.Stop()
	log.Infof("GRPC Server was shutdown")
}
