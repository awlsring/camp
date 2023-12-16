package grpc

type ServerOpt func(*CampdGrpcServer)

func WithNetwork(network string) ServerOpt {
	return func(s *CampdGrpcServer) {
		s.network = network
	}
}

func WithAddress(address string) ServerOpt {
	return func(s *CampdGrpcServer) {
		s.address = address
	}
}
