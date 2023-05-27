package grpc

import (
	"context"
	"fmt"

	"my-microservice/protos"
)

type GreeterService struct {
	protos.UnimplementedGreeterServer
}

func (g *GreeterService) SayHello(ctx context.Context, request *protos.HelloRequest) (*protos.HelloReply, error) {
	return &protos.HelloReply{
		Message: fmt.Sprintf("Hello there, %s", request.Name),
	}, nil
}
