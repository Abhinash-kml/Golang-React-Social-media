package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/utils"
)

func CookieBasedJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("Cookie JWT failed. No cookie found.")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		token := cookie.Value
		verifiedToken, err := utils.VerifyJWT(token)
		if err != nil {
			fmt.Printf("Token verification failed: %v\\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Print information about the verified token
		fmt.Printf("Token verfied successfully. Claims: %+v\\n", verifiedToken.Claims)

		next.ServeHTTP(w, r)
	})
}

func HeaderBasedJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		token := strings.Split(bearerToken, " ")[1]
		_, err := utils.VerifyJWT(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		jwt.
			w.WriteHeader(http.StatusOK)
		next.ServeHTTP(w, r)
	})
}
