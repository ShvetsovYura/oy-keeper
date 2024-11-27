package server

type GRPCServer struct {
	server  grpc.server
	address string
}

func NewGRPCServer() {
	return &GRPCServer{
		grpcServer: *grpc.NewServer()
	}
}
func (s *GRPCServer) RegisterHandlers(targetStorage handlers.Storage, opt *Options) {
	pb.RegisterMetricsServer(
		&s.grpcServer,
		handlers.NewMetricServer(targetStorage),
	)
	s.addr = opt.EndpointAddr
}

func (s *GRPCServer) StartListen() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(listen)

}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	s.grpcServer.Stop()
	return nil
}

func (s *KServer) UploadFile