package signalling

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// ServeWS ...
func ServeWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["id"]

	//TODO(james): Check user authentication to join class

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading websocket: %s", err)
	}

	// Create a new classroom if one for this code does not already exist
	classrooms.mu.Lock()
	if _, ok := classrooms.classes[classID]; !ok {
		log.Printf("Creating Classroom: %s", classID)
		classrooms.classes[classID] = &Class{
			Code:    classID,
			Members: []*websocket.Conn{},
		}
	}

	log.Infof("Connecting %s to: %s", r.RemoteAddr, classID)

	// Add connection to class
	classrooms.classes[classID].Members = append(classrooms.classes[classID].Members, ws)
	classrooms.mu.Unlock()
	wsHandler(ws, classrooms.classes[classID])
}
