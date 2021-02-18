package routes

import (
	"errors"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//TODO(james): Origin checking
	CheckOrigin: func(r *http.Request) bool { return true },
}

func InjectSignallingRoutes(subrouter *mux.Router) {
	//TODO(james): Authentication
	// Connect to WebSocket
	subrouter.HandleFunc("/ws/{classroomId}", joinClassroom)
	// Turn Server Credentials
	subrouter.HandleFunc("/credentials", credentials)
}

func joinClassroom(w http.ResponseWriter, r *http.Request) {
	id, err := getClassroomId(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
	}

	services.SignallingAddToClassroom(ws, id)
}

func credentials(w http.ResponseWriter, r *http.Request) {
	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	id := authContext.Account.ID
	WriteBody(w, r, services.GenerateTURNCredentials(id))
}

func getClassroomId(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	if val, ok := vars["classroomId"]; ok {
		return val, nil
	}
	return "", errors.New("No Classroom ID Specified")
}
