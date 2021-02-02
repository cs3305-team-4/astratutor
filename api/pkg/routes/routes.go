package routes

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func GetHandler() http.Handler {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	InjectAccountsRoutes(r.PathPrefix("/accounts").Subrouter())
	InjectAuthRoutes(r.PathPrefix("/auth").Subrouter())
	InjectLessonsRoutes(r.PathPrefix("/lessons").Subrouter())
	InjectStudentsRoutes(r.PathPrefix("/students").Subrouter())
	InjectSubjectsRoutes(r.PathPrefix("/subjects").Subrouter())
	InjectTutorsRoutes(r.PathPrefix("/tutors").Subrouter())

	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*",
		},
	}).Handler(r)
}

const (
	sqlDuplicate = "(SQLSTATE 23505)"
)

type returnError struct {
	Error          string `json:"error"`
	Type           string `json:"type"`
	DuplicateField string `json:"duplicate_field,omitempty"`
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.WithContext(r.Context()).WithError(err).Error("Error parsing body")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	var outErr returnError

	// Special error cases
	switch {
	case strings.Contains(err.Error(), sqlDuplicate):
		re, e := regexp.Compile("_([a-z]+)_key")
		if e != nil {
			log.WithContext(r.Context()).WithError(err).Error("Error parsing body")
		}
		match := re.FindStringSubmatch(err.Error())
		if len(match) > 1 {
			outErr = returnError{
				Error:          "Duplicate key found",
				DuplicateField: match[1],
				Type:           "duplicate",
			}
		}
	default:
		outErr = returnError{
			Error: err.Error(),
			Type:  "standard",
		}
	}

	if err := json.NewEncoder(w).Encode(&outErr); err != nil {
		log.WithContext(r.Context()).WithError(err).Error("Could not return error to user")
		return
	}
}
