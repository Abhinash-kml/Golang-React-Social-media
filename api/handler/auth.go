package handler

import (
	"database/sql"
	"net/http"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/models"
	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/signup", SignupHandler).Methods("POST")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	passwordFromDb, err := models.GetPasswordOfUserWithEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Other type of error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored password in db with the hashed password currently created from request
	if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(passwordFromDb)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	// Password matched, create token and proceed
	token, err := utils.CreateJWT(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create cookie
	//loggedInUser = email
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, cookie)                     // Set cookie
	http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = InsertNewUserIntoDatabase(email, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
