package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterAndStartApiServer() {

	router := mux.NewRouter()
	usersRouter := router.PathPrefix("/api/users").Subrouter()
	postRouter := router.PathPrefix("/api/posts").Subrouter()
	messageRouter := router.PathPrefix("/api/message").Subrouter()
	mediaRouter := router.PathPrefix("/api/media").Subrouter()

	RegisterUserRoutes(usersRouter)
	RegisterPostRoutes(postRouter)
	RegisterMessageRoutes(messageRouter)
	RegisterMediaRoutes(mediaRouter)

	http.ListenAndServe(":8000", router)
}
