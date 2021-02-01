package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func GetHandler() http.Handler {
	r := mux.NewRouter()

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
