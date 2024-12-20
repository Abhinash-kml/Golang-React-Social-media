package server

import (
	"net/http"

	"github.com/Abhinash-kml/Golang-React-Social-media/internal/server/api/handler"
	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	repo   *db.Postgres
	router *mux.Router
}

func NewServer() *Server {
	newlogger, _ := zap.NewProduction()

	return &Server{
		logger: newlogger,
		repo:   &db.Postgres{},
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() {
	s.InitializeDatabaseConnection()
	s.SetupRoutes()
}

func (s *Server) Stop() {
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
