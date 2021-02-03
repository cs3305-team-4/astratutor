package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func InjectLessonsRoutes(subrouter *mux.Router) {
	subrouter.PathPrefix("/{uuid}").Methods("GET").HandlerFunc(handleLessonGet)

}

func handleLessonGet(W http.ResponseWriter, r *http.Request) {

}
