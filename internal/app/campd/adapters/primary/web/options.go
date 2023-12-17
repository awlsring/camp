package web

type ServerOpt func(*Server)

func WithAddress(address string) ServerOpt {
	return func(s *Server) {
		s.address = address
	}
}
