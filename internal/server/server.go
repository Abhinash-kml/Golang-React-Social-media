package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Abhinash-kml/Golang-React-Social-media/internal/server/api/handler"
	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	logger     *zap.Logger
	repository db.Repository
	router     *mux.Router
	httpserver *http.Server
}

func NewServer() *Server {
	newlogger, _ := zap.NewProduction()
	muxRouter := mux.NewRouter()

	return &Server{
		logger:     newlogger,
		repository: &db.Postgres{},
		router:     muxRouter,
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
	s.repository.Disconnect()
}

func (s *Server) InitializeDatabaseConnection() {
	s.repository.Connect()
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	handler.HandleLogin(s.repository, w, r)
}

func (s *Server) HandleSignup(w http.ResponseWriter, r *http.Request) {
	handler.HandleSignup(s.logger, s.repository, w, r)
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

// Tested - OK
func (s *Server) GetUserWithAttribute(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	attributeType := queryParams.Get("attribute_type")
	attribute := queryParams.Get("attribute")

	validAttributes := []string{
		"userid",
		"name",
		"email",
		"password",
		"dob",
		"created_at",
		"modified_at",
		"last_login",
		"country",
		"state",
		"city",
		"avatar_url",
		"ban_level",
		"ban_duration",
	}

	valid := false
	for _, val := range validAttributes {
		if attributeType == val {
			valid = true
			break
		}
	}

	if !valid {
		http.Error(w, "invalid attribute", http.StatusInternalServerError)
		return
	}

	users, err := s.repository.GetUsersWithAttribute(context.Background(), attributeType, attribute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(users)
	w.WriteHeader(http.StatusOK)
}

// Tested - OK
func (s *Server) GetPostsOfUserid(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userid := queryParams.Get("userid")
	uuid, err := uuid.Parse(userid)
	if err != nil {
		http.Error(w, "Failed to parse provided uuid", http.StatusInternalServerError)
		return
	}

	posts, err := s.repository.GetPostsOfUser(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
	w.WriteHeader(http.StatusOK)
}

// Tested - OK
func (s *Server) AddNewUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	dob := r.FormValue("dob")
	country := r.FormValue("country")
	state := r.FormValue("state")
	city := r.FormValue("city")
	avatarurl := r.FormValue("avatarurl")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		http.Error(w, "Hashing password failed", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.InsertUser(context.Background(), name, email, string(hashedPassword), dob, country, state, city, avatarurl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Testing Required
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	newName := r.FormValue("name")
	newEmail := r.FormValue("email")
	newCountry := r.FormValue("country")
	newState := r.FormValue("state")

	userid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.UpdateUserWithId(context.Background(), userid, newName, newEmail, newCountry, newState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Testing required
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.DeleteUserWithId(context.Background(), userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Tested - OK
func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.repository.GetAllUsers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetCommentsOfPostId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postid := vars["id"]

	uuid, err := uuid.Parse(postid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	comments, err := s.repository.GetCommentsOfPost(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	comments, err := s.repository.GetAllPosts(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetPostWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postid := vars["id"]

	uuid, err := uuid.Parse(postid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	post, err := s.repository.GetPostWithId(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) UpdatePostWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postid := vars["id"]
	newTitle := r.FormValue("title")
	newBody := r.FormValue("body")
	newHashtag := r.FormValue("hashtag")

	uuid, err := uuid.Parse(postid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.UpdatePostWithId(context.Background(), uuid, newTitle, newBody, newHashtag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal Query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeletePostWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postid := vars["id"]

	uuid, err := uuid.Parse(postid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.DeletePostWithId(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal Query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetPostsOfUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["id"]

	uuid, err := uuid.Parse(userid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	posts, err := s.repository.GetPostsOfUser(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) AddPostOfUserWithId(w http.ResponseWriter, r *http.Request) {
	// Extract form values
	userid := r.FormValue("userid")
	title := r.FormValue("title")
	body := r.FormValue("body")
	mediaUrl := r.FormValue("media_url")
	hashtag := r.FormValue("hashtag")

	// Parse the UUID from the string
	uuid, err := uuid.Parse(userid)
	if err != nil {
		http.Error(w, "Invalid UUID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Insert post into the repository
	ok, err := s.repository.InsertPost(context.Background(), uuid, title, body, mediaUrl, hashtag)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal Query operation failed", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func (s *Server) AddCommentToPostWithId(w http.ResponseWriter, r *http.Request) {
	postid := r.FormValue("postid")
	body := r.FormValue("body")

	uuid, err := uuid.Parse(postid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.AddCommentToPostId(context.Background(), uuid, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal query operation failed", http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateCommentWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentid := vars["id"]
	newBody := r.FormValue("body")

	uuid, err := uuid.Parse(commentid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.UpdateCommentWithId(context.Background(), uuid, newBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal Query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteCommentWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentid := vars["id"]

	uuid, err := uuid.Parse(commentid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	ok, err := s.repository.DeleteCommentWithId(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Internal Query operation failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetAllComments(w http.ResponseWriter, r *http.Request) {

	comments, err := s.repository.GetAllComments(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetCommentWithId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentid := vars["id"]

	uuid, err := uuid.Parse(commentid)
	if err != nil {
		http.Error(w, "Error parsing uuid", http.StatusInternalServerError)
		return
	}

	comment, err := s.repository.GetCommentWithId(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetRepo() db.Repository {
	return s.repository
}
