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
	if _, ok := classes[classID]; !ok {
		log.Infof("Creating Classroom: %s", classID)
		classes[classID] = &Classroom{
			Code:    classID,
			Members: []*websocket.Conn{},
		}
	}

	log.Infof("Connecting %s to: %s", r.RemoteAddr, classID)

	// Add connection to class
	classes[classID].Members = append(classes[classID].Members, ws)
	wsHandler(ws, classes[classID])
}
