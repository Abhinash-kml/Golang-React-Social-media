package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Abhinash-kml/Golang-React-Social-media/internal/server/api/handler"
	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	logger     *zap.Logger
	repo       *db.Postgres
	router     *mux.Router
	httpserver *http.Server
}

func NewServer() *Server {
	newlogger, _ := zap.NewProduction()
	muxRouter := mux.NewRouter()

	return &Server{
		logger: newlogger,
		repo:   &db.Postgres{},
		router: muxRouter,
		httpserver: &http.Server{
			Addr:         ":8000",
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      muxRouter,
		},
	}
}

func (s *Server) Start() {
	s.InitializeDatabaseConnection()
	s.SetupRoutes()
	s.ServeAPI()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	s.httpserver.Shutdown(ctx)
	s.repo.Disconnect()
}

func (s *Server) InitializeDatabaseConnection() {
	s.repo.Connect()
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	handler.HandleLogin(s.logger, s.repo, w, r)
}

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	handler.HandleSignup(s.logger, s.repo, w, r)
}

func (s *Server) PrivateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Passed cookie jwt auth")

	json.NewEncoder(w).Encode("Auth: Passed")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) ServeAPI() {
	go func() {
		if err := s.httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start API server. Error:", err)
			s.Stop()
		}
	}()

	fmt.Println("Listening on localhost:8000.")
}

func (s *Server) GetCommentsOfPostId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetUserWithAttribute(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	attributeType := queryParams.Get("attribute_type")
	attribute := queryParams.Get("attribute")

	switch attributeType {
	case "userid":
		{
			fmt.Println("attribute is userid")
			user, err := s.repo.GetUserWithID(s.logger, attribute)
			if err != nil {
				s.logger.Error("Error getting user with userid from database",
					zap.String("Error", err.Error()))

				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(user)
			w.WriteHeader(http.StatusOK)
		}
	case "name":
		{
			fmt.Println("attribute is name")
			user, err := s.repo.GetUserWithName(s.logger, attribute)
			if err != nil {
				s.logger.Error("Error getting user with name from database",
					zap.String("Error", err.Error()))

				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(user)
			w.WriteHeader(http.StatusOK)
		}
	case "email":
		fmt.Println("attribute is email")
		user, err := s.repo.GetUserWithEmail(s.logger, attribute)
		if err != nil {
			s.logger.Error("Error getting user with email from database",
				zap.String("Error", err.Error()))

			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) GetPostsOfUserid(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) AddNewUser(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) AddCommentToPostWithId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) AddNewCommentToPostWithId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UpdateCommentOfPostWithId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) DeleteCommentOfPostWithId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) DeletePostOfUserWithPostId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) UpdatePostOfUserWithPostId(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) AddPostOfUser(w http.ResponseWriter, r *http.Request) {

}
