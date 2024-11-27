package server

import (
	"context"
	"net"

	"github.com/ShvetsovYura/oy-keeper/internal/server/handlers"
	pb "github.com/ShvetsovYura/oy-keeper/proto"
	"google.golang.org/grpc"
)

type Server struct {
	server grpc.Server
	addr   string
}

func New() *Server {
	return &Server{
		server: *grpc.NewServer(),
		addr:   ":9001",
	}
}

func (s *Server) Run() error {
	pb.RegisterFileServiceServer(&s.server, handlers.NewFileUploadHandler())

	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.server.Serve(listen)
}

func (s *Server) Shutdown(ctx context.Context) {
	s.server.GracefulStop()
}
