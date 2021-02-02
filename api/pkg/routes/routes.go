package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

func GetHandler() http.Handler {
	services.SetCustomValidators()
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

// ParseBody inplace. Returns false if error has been written.
func ParseBody(w http.ResponseWriter, r *http.Request, i interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return false
	}
	if err := validator.Validate(i); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return false
	}
	return true
}

const (
	sqlDuplicate = "(SQLSTATE 23505)"
)

type returnDetail struct {
	Loc []string `json:"loc"`
	Msg string   `json:"msg"`
}
type returnError struct {
	Detail returnDetail `json:"detail"`
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.WithContext(r.Context()).WithError(err).Error("Error parsing body")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	// Error cases
	switch {
	case strings.Contains(err.Error(), sqlDuplicate):
		re, e := regexp.Compile("_([a-z]+)_key")
		if e != nil {
			log.WithContext(r.Context()).WithError(err).Error("Error parsing body")
		}
		match := re.FindStringSubmatch(err.Error())
		if len(match) > 1 {
			err = fmt.Errorf("%s already exists", match[1])
		}
	}
	outErr := returnError{
		Detail: returnDetail{
			Msg: err.Error(),
		},
	}

	if err := json.NewEncoder(w).Encode(&outErr); err != nil {
		log.WithContext(r.Context()).WithError(err).Error("Could not return error to user")
		return
	}
}
