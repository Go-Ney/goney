package transport

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server   *grpc.Server
	listener net.Listener
	port     string
}

type GrpcService interface {
	RegisterWithServer(*grpc.Server)
}

func NewGrpcServer(port string) *GrpcServer {
	return &GrpcServer{
		server: grpc.NewServer(),
		port:   port,
	}
}

func (s *GrpcServer) RegisterService(service GrpcService) {
	service.RegisterWithServer(s.server)
}

func (s *GrpcServer) EnableReflection() {
	reflection.Register(s.server)
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}
	s.listener = lis

	return s.server.Serve(lis)
}

func (s *GrpcServer) Stop() {
	s.server.GracefulStop()
}

func (s *GrpcServer) StopNow() {
	s.server.Stop()
}

type BaseGrpcService struct{}

func (b *BaseGrpcService) HandleError(ctx context.Context, err error) error {
	return err
}