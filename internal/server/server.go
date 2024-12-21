package server

import (
	"database/sql"
	"net/http"

	"github.com/Abhinash-kml/Golang-React-Social-media/internal/server/api/handler"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	Logger     *zap.Logger
	Connection *sql.DB
	Router     *mux.Router
}

func NewServer() *Server {
	logger, _ := zap.NewProduction()

	return &Server{
		Logger:     logger,
		Connection: nil,
		Router:     mux.NewRouter(),
	}
}

func (s *Server) Start() {
	s.SetupRoutes()
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	handler.HandleLogin(s.Logger, w, r)
}

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	handler.HandleSignup(s.Logger, w, r)
}
