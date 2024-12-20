package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = []byte("your-secret-key")
var loggedInUSer string

func RegisterAndStartApiServer() {
	router := mux.NewRouter()
	usersRouter := router.PathPrefix("/api/users").Subrouter()

	RegisterUserRoutes(usersRouter)

	http.ListenAndServe(":8000", router)
}

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/signup", SignupHandler).Methods("POST")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	passwordFromDb, err := GetPasswordOfUserWithEmail(email)
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
	if err = bcrypt.CompareHashAndPassword(hashedPassword, passwordFromDb); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	// Password matched, create token and proceed
	token, err := CreateJWT(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create cookie
	loggedInUSer = email
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

func CreateJWT(username string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "social-media",                   // Issuer
		"aud": "coder",                          // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	fmt.Printf("Token claims added: %+v", claims)

	tokenString, err := claims.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		token := cookie.Value
		verifiedToken, err := VerifyJWT(token)
		if err != nil {
			fmt.Printf("Token verification failed: %v\\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Print information about the verified token
		fmt.Printf("Token verfied successfully. Claims: %+v\\n", verifiedToken.Claims)
	})
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}
