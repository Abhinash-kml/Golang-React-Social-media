package server

func (s *Server) SetupRoutes() {
	accountRouter := s.router.PathPrefix("/api").Subrouter()
	accountRouter.HandleFunc("/login", s.HandleLogin).Methods("POST")
	accountRouter.HandleFunc("/signup", s.HandleSignup).Methods("POST")
}
