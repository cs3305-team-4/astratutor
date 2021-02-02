package routes

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// // Define our struct
// type authenticationMiddleware struct {
// 	tokenUsers map[string]string
// }

// type authJwt struct {

// }

// // Initialize it somewhere
// func (amw *authenticationMiddleware) Populate() {
// 	amw.tokenUsers["00000000"] = "user0"
// 	amw.tokenUsers["aaaaaaaa"] = "userA"
// 	amw.tokenUsers["05f717e5"] = "randomUser"
// 	amw.tokenUsers["deadbeef"] = "user0"
// }

// // Middleware function, which will be called for each request
// func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("X-Session-Token")

// 		if user, found := amw.tokenUsers[token]; found {
// 			// We found the token in our map
// 			log.Printf("Authenticated user %s\n", user)
// 			next.ServeHTTP(w, r)
// 		} else {
// 			http.Error(w, "Forbidden", http.StatusForbidden)
// 		}
// 	})
// }

// func Example_authenticationMiddleware() {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// Do something here
// 	})
// 	amw := authenticationMiddleware{make(map[string]string)}
// 	amw.Populate()
// 	r.Use(amw.Middleware)
// }

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)) {

// 	}
// }

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithContext(r.Context()).Infof("Incoming request to %s", r.URL.String())
		next.ServeHTTP(w, r)
	})
}
