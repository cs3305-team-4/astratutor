package routes

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func authAccount() func(next http.Handler) http.Handler {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error {
		id, err := getUUID(r, "uuid")
		if err != nil {
			return err
		}

		if ac.Account.ID != id {
			return errors.New("cannot operate on a resource you do not own")
		}

		return nil
	}, true)
}

func authSetCtx() func(next http.Handler) http.Handler {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error {
		return nil
	}, false)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"host":           r.Host,
			"slug":           r.URL.Path,
			"method":         r.Method,
			"user_agent":     r.Header.Get("User-Agent"),
			"source_address": r.RemoteAddr,
		}).WithContext(r.Context()).Info("Incoming request")
		next.ServeHTTP(w, r)
	})
}
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}
