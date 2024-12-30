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
	private := s.router.PathPrefix("/api/private").Subrouter()
	private.Use(middleware.LoggingMiddleware)
	// private.Use(middleware.CookieBasedJWTAuth)
	// private.Use(middleware.RateLimit)

	// User Routes
	public.HandleFunc("/users", s.GetAllUsers).Methods("GET")               // working
	public.HandleFunc("/users", s.AddNewUser).Methods("POST")               // working
	public.HandleFunc("/users/{id}", s.GetUserWithAttribute).Methods("GET") // works :: refinement and flexible query design required
	public.HandleFunc("/users/{id}", s.UpdateUser).Methods("PUT")           // working
	public.HandleFunc("/users/{id}", s.DeleteUser).Methods("DELETE")        // working

	// Post Routes
	public.HandleFunc("/posts", s.GetAllPosts).Methods("GET")                     // working
	public.HandleFunc("/posts/{id}", s.GetPostWithId).Methods("GET")              // working
	public.HandleFunc("/posts/{id}", s.UpdatePostWithId).Methods("PUT")           // working
	public.HandleFunc("/posts/{id}", s.DeletePostWithId).Methods("DELETE")        // working
	public.HandleFunc("/users/{id}/posts", s.GetPostsOfUserId).Methods("GET")     // working
	public.HandleFunc("/users/{id}/posts", s.AddPostOfUserWithId).Methods("POST") // working

	// Comment Routes
	public.HandleFunc("/comments", s.GetAllComments).Methods("GET")                     // working
	public.HandleFunc("/comments/{id}", s.GetCommentWithId).Methods("GET")              // working
	public.HandleFunc("/comments/{id}", s.UpdateCommentWithId).Methods("PUT")           // working
	public.HandleFunc("/comments/{id}", s.DeleteCommentWithId).Methods("DELETE")        // working
	public.HandleFunc("/posts/{id}/comments", s.GetCommentsOfPostId).Methods("GET")     // testing required
	public.HandleFunc("/posts/{id}/comments", s.AddCommentToPostWithId).Methods("POST") // testing required
}
