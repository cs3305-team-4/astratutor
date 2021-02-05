package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/validator.v2"
)

func GetHandler() http.Handler {
	services.SetCustomValidators()
	r := mux.NewRouter()

	r.Use(
		loggingMiddleware,
		jsonMiddleware,
	)

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
		restError(w, r, err, http.StatusBadRequest)
		return false
	}
	if err := validator.Validate(i); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return false
	}
	return true
}

// getUUID can parse a UUID from the router variables
// if param is nil, the default variable used "uuid"
func getUUID(r *http.Request, param string) (uuid.UUID, error) {
	para := "uuid"
	if param != "" {
		para = param
	}

	vars := mux.Vars(r)
	val, ok := vars[para]
	if !ok {
		return uuid.UUID{}, errors.New("no uuid found")
	}
	return uuid.Parse(val)
}
