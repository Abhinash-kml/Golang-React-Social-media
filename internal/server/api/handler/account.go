package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(repository db.Repository, w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	passwordFromDb, err := repository.GetPasswordOfUserWithEmail(context.Background(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("No sql row.")
			return
		}

		// Other type of error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored password in db with the hashed password currently created from request
	if err = bcrypt.CompareHashAndPassword([]byte(passwordFromDb), []byte(password)); err != nil {
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

func HandleSignup(logger *zap.Logger, repository db.Repository, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	logger.Info("Recieved data",
		zap.String("name", name),
		zap.String("email", email),
		zap.String("password", password))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ok, err := repository.InsertUser(context.Background(), name, email, string(hashedPassword), "2001-09-13", "India", "Bengal", "Kolkata", "www.imgur.com")
	if !ok || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
