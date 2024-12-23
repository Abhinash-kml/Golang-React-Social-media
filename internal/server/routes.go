package server

import "github.com/Abhinash-kml/Golang-React-Social-media/internal/server/api/middleware"

func (s *Server) SetupRoutes() {
	account := s.router.PathPrefix("/account").Subrouter()
	account.Use(middleware.LoggingMiddleware)
	// account.Use(middleware.RateLimit)

	account.HandleFunc("/login", s.HandleLogin).Methods("POST")
	account.HandleFunc("/signup", s.HandleSignup).Methods("POST")

	public := s.router.PathPrefix("/api/public").Subrouter()
	public.Use(middleware.LoggingMiddleware)
	// public.Use(middleware.RateLimit)
	// public.HandleFunc("/user", s.GetUserWithName).Methods("GET")
	// public.HandleFunc("/user", s.AddNewUser).Methods("POST")
	// public.HandleFunc("/users", s.GetAllUsers).Methods("GET")
	// public.HandleFunc("/users", s.AddUsers).Methods("POST")

	private := s.router.PathPrefix("/api/private").Subrouter()
	private.Use(middleware.LoggingMiddleware)
	private.Use(middleware.PerformCookieBasedJWTAuth)
	private.HandleFunc("/user", s.PrivateHandler).Methods("GET")
	// private.Use(middleware.RateLimit)

	// userRouter := s.router.PathPrefix("/api/user").Subrouter()
	// userRouter.HandleFunc("", s.HandleGetUser).Methods("GET")
	// userRouter.HandleFunc("", s.HandleSetUser).Methods("SET")

	// postRouter := s.router.PathPrefix("/api/post").Subrouter()

	// commentRouter := s.router.PathPrefix("/api/comment").Subrouter()

	// mediaRouter := s.router.PathPrefix("/api/media").Subrouter()
}
