package signalling

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(conn *websocket.Conn, class *Classroom) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Disconnected
			log.Error(err)
			for i, member := range class.Members {
				if member == conn {
					class.Members = append(class.Members[:i], class.Members[i+1:]...)
					// Remove Class when all members have disconnected
					if len(class.Members) == 0 {
						log.Infof("Deleting Classroom: %s", class.Code)
						delete(classes, class.Code)
					}
					return
				}
			}
		}

		// Send received message to all other members of the class
		for _, c := range class.Members {
			if conn == c {
				continue
			}
			c.WriteMessage(messageType, message)
		}
	}
}
