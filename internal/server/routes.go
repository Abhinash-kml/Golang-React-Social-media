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
	//private.Use(middleware.CookieBasedJWTAuth)
	// private.Use(middleware.RateLimit)

	public.HandleFunc("/user", s.GetUserWithAttribute).Methods("GET")
	public.HandleFunc("/comments/{postid}", s.GetCommentsOfPostId).Methods("GET")
	public.HandleFunc("/posts/{userid}", s.GetPostsOfUserid).Methods("GET")

	private.HandleFunc("/user", s.AddNewUser).Methods("POST")
	private.HandleFunc("/users", s.GetAllUsers).Methods("GET")
	private.HandleFunc("/comment", s.AddNewCommentToPostWithId).Methods("POST")
	private.HandleFunc("/comment/{postid}", s.AddCommentToPostWithId).Methods("POST")
	private.HandleFunc("/comment/{postid}/{id}", s.UpdateCommentOfPostWithId).Methods("PUT")
	private.HandleFunc("/comment/{postid}/{id}", s.DeleteCommentOfPostWithId).Methods("DELETE")
	private.HandleFunc("/posts/{userid}/{postid}", s.DeletePostOfUserWithPostId).Methods("DELETE")
	private.HandleFunc("/posts/{userid}/{postid}", s.UpdatePostOfUserWithPostId).Methods("PUT")
	private.HandleFunc("/posts", s.AddPostOfUserid).Methods("POST")

	// userRouter := s.router.PathPrefix("/api/user").Subrouter()
	// userRouter.HandleFunc("", s.HandleGetUser).Methods("GET")
	// userRouter.HandleFunc("", s.HandleSetUser).Methods("SET")

	// postRouter := s.router.PathPrefix("/api/post").Subrouter()

	// commentRouter := s.router.PathPrefix("/api/comment").Subrouter()

	// mediaRouter := s.router.PathPrefix("/api/media").Subrouter()
}
