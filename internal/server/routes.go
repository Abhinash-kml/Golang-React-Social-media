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
	public.HandleFunc("/users", s.GetAllUsers).Methods("GET")               // Query params: ?limit=10&offset=0
	public.HandleFunc("/users/{id}", s.GetUserWithAttribute).Methods("GET") // Path param: {id}
	public.HandleFunc("/users", s.AddNewUser).Methods("POST")               // Add new user
	public.HandleFunc("/users/{id}", s.UpdateUser).Methods("PUT")           // Update existing user
	public.HandleFunc("/users/{id}", s.DeleteUser).Methods("DELETE")        // Delete existing user

	// Post Routes
	public.HandleFunc("/posts", s.GetAllPosts).Methods("GET")        // Query params: ?user_id=123&hashtag=sports
	public.HandleFunc("/posts/{id}", s.GetPostWithId).Methods("GET") // Path param: {id}
	public.HandleFunc("/posts/{id}", s.UpdatePostWithId).Methods("PUT")
	public.HandleFunc("/posts/{id}", s.DeletePostWithId).Methods("DELETE")
	public.HandleFunc("/users/{id}/posts", s.GetPostsOfUserId).Methods("GET")     // Query params: ?limit=10&offset=0
	public.HandleFunc("/users/{id}/posts", s.AddPostOfUserWithId).Methods("POST") // Query params: ?limit=10&offset=0

	// Comment Routes
	public.HandleFunc("/comments", s.GetAllComments).Methods("GET")        // Query params: ?post_id=123
	public.HandleFunc("/comments/{id}", s.GetCommentWithId).Methods("GET") // Path param: {id}
	public.HandleFunc("/comments/{id}", s.UpdateCommentWithId).Methods("PUT")
	public.HandleFunc("/comments/{id}", s.DeleteCommentWithId).Methods("DELETE")
	public.HandleFunc("/posts/{id}/comments", s.GetCommentsOfPostId).Methods("GET")     // Query params: ?limit=10&offset=0
	public.HandleFunc("/posts/{id}/comments", s.AddCommentToPostWithId).Methods("POST") // Path param: {id}
}
