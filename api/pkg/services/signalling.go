package services

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var classrooms = Classrooms{
	classrooms: map[string]*Classroom{},
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Classrooms ...
type Classrooms struct {
	classrooms map[string]*Classroom
	mu         sync.Mutex
}

// Classroom ...
type Classroom struct {
	Code    string
	Members []*websocket.Conn
	mu      sync.Mutex
}

func SignallingAddToClassroom(ws *websocket.Conn, classroomId string) {
	classrooms.mu.Lock()
	if _, ok := classrooms.classrooms[classroomId]; !ok {
		log.Infof("Creating Classroom: %s", classroomId)
		classrooms.classrooms[classroomId] = &Classroom{
			Code:    classroomId,
			Members: []*websocket.Conn{},
		}
	}
	log.Infof("%s connecting to %s", ws.RemoteAddr(), classroomId)
	classroom := classrooms.classrooms[classroomId]
	classroom.Members = append(classroom.Members, ws)
	classrooms.mu.Unlock()
	messageHandler(ws, classroom)
}

func messageHandler(conn *websocket.Conn, class *Classroom) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Disconnected
			log.Error(err)
			class.mu.Lock()
			for i, member := range class.Members {
				if member == conn {
					class.Members = append(class.Members[:i], class.Members[i+1:]...)
					// Remove Class when all members have disconnected
					if len(class.Members) == 0 {
						log.Infof("Deleting Classroom: %s", class.Code)
						classrooms.mu.Lock()
						delete(classrooms.classrooms, class.Code)
						classrooms.mu.Unlock()
					}
					class.mu.Unlock()
					return
				}
			}
		}

		// Send received message to all other members of the class
		class.mu.Lock()
		for _, c := range class.Members {
			if conn == c {
				continue
			}
			c.WriteMessage(messageType, message)
		}
		class.mu.Unlock()
	}
}

func GenerateTURNCredentials(id uuid.UUID) Credentials {
	// Valid for 1.5 Hours
	timestamp := time.Now().Unix() + (60 * 90)
	username := fmt.Sprintf("%d:%s", timestamp, id)

	key := viper.GetString("signalling_secret")

	hmac := hmac.New(crypto.SHA1.New, []byte(key))
	hmac.Write([]byte(username))
	password := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	return Credentials{
		Username: username,
		Password: password,
	}
}
