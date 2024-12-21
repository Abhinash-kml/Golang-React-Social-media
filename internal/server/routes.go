package server

func (s *Server) SetupRoutes() {
	accountRouter := s.Router.PathPrefix("/api").Subrouter()
	accountRouter.HandleFunc("/login", s.HandleLogin).Methods("POST")
	accountRouter.HandleFunc("/signup", s.HandleSignup).Methods("POST")
}
