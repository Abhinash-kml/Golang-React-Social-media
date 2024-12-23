package server

import "github.com/Abhinash-kml/Golang-React-Social-media/internal/server/api/middleware"

func (s *Server) SetupRoutes() {
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.LoggingMiddleware)
	apiRouter.HandleFunc("/account/login", s.HandleLogin).Methods("POST")
	apiRouter.HandleFunc("/account/signup", s.HandleSignup).Methods("POST")

	// userRouter := s.router.PathPrefix("/api/user").Subrouter()
	// userRouter.HandleFunc("", s.HandleGetUser).Methods("GET")
	// userRouter.HandleFunc("", s.HandleSetUser).Methods("SET")

	// postRouter := s.router.PathPrefix("/api/post").Subrouter()

	// commentRouter := s.router.PathPrefix("/api/comment").Subrouter()

	// mediaRouter := s.router.PathPrefix("/api/media").Subrouter()
}
