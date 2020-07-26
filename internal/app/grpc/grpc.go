package grpc

import (
	"context"
	"fmt"
	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/chunganhbk/gin-go/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)
// ServerOptions are shared by all grpc servers
var ServerOptions = []grpc.ServerOption{
	// XXX: this is done to prevent routers from cleaning up our connections (e.g aws load balances..)
	// TODO: these parameters work for now but we might need to revisit or add them as configuration
	// TODO: Configure maxconns, maxconcurrentcons ..
	grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Minute * 120,
		MaxConnectionAge:      time.Minute * 180,
		MaxConnectionAgeGrace: time.Minute * 10,
		Time:                  time.Minute,
		Timeout:               time.Minute * 3,
	}),
}
func NewRpcServer() *RpcServer {
	port := config.C.GrpcServer.Port | 7777
	return &RpcServer{
		Port:       port,
		Server: grpc.NewServer(ServerOptions...),
	}
}
// is a very basic grpc server
type RpcServer struct {
	Port       int
	Server *grpc.Server
}
func (s *RpcServer) Start(ctx context.Context) {
	logger.Infof(ctx,"starting new grpc server")
	go s.startInternal(ctx)
}


// Blocking, should be called in a goroutine
func (s *RpcServer) startInternal(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		logger.Errorf(ctx, "error listening on port", err)
		return
	}

	// SubscribeOnNewConnections reflection service on gRPC server
	reflection.Register(s.Server)

	// start serving - this blocks until err or server is stopped
	logger.Infof(ctx, "starting new grpc server on port %d", s.Port)
	if err := s.Server.Serve(lis); err != nil {
		logger.Errorf(ctx, "error stopping grpc server: %v", err)
	}
}

// Close stops the server
func (s *RpcServer) Close(ctx context.Context) {
	logger.Infof(ctx, "Stopping new grpc server...")
	s.Server.Stop()
}
// Ser
