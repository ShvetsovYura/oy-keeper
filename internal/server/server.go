package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
	"github.com/ShvetsovYura/oykeeper/internal/server/service"
	"github.com/ShvetsovYura/oykeeper/internal/server/store"
	pb "github.com/ShvetsovYura/oykeeper/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedFileServiceServer
	server grpc.Server
	addr   string
}

func New() *Server {
	return &Server{
		server: *grpc.NewServer(),
		addr:   ":9091",
	}
}

func (s *Server) Run(ctx context.Context) error {

	dbDSN := "postgres://lbman:lbman@localhost:5432/lb1"
	conn, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	err = conn.Ping(ctx)
	if err != nil {
		panic(err.Error())
	}
	recordStore := store.NewRecordStore(conn)
	pb.RegisterFileServiceServer(&s.server, service.NewFileService("hoho"))
	pb.RegisterRecordServiceServer(&s.server, service.NewRecordService(recordStore))

	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.Shutdown()
	}()
	if err := s.server.Serve(listen); err != nil {
		return fmt.Errorf("error: %w", err)
	}
	wg.Wait()
	return nil
}

func (s *Server) Shutdown() {
	logger.Log.Info("graceful shudown")
	s.server.GracefulStop()
}
