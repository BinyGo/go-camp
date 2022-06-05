package server

import (
	"context"
	"net"

	"github.com/blog/internal/service"
	"github.com/blog/pkg/app"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	listener net.Listener
	server   *grpc.Server
}

func (g *GrpcServer) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		g.server.Stop()
	}()
	return g.server.Serve(g.listener)
}

func NewGrpcServer(service *service.BlogService) app.GrpcServer {
	server := new(GrpcServer)
	lis, err := net.Listen("tcp", ":8999")
	server.listener = lis
	if err != nil {
		panic(err.Error())
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	//加载api proto GRPCServer
	//v1.RegisterBlogServer(grpcServer, service)
	server.server = grpcServer
	return server
}
