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
