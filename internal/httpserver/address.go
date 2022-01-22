package httpserver

// GetAddress obtains the address the HTTP server is listening on.
func (s *Server) GetAddress() (address string) {
	<-s.addressSet
	return s.address
}
