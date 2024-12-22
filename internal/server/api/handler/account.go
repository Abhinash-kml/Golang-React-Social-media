package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(logger *zap.Logger, repo *db.Postgres, w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	passwordFromDb, err := repo.GetPasswordOfUserWithEmail(logger, email)
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

func HandleSignup(logger *zap.Logger, repo *db.Postgres, w http.ResponseWriter, r *http.Request) {
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

	err = repo.InsertNewUserIntoDatabase(logger, name, email, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
