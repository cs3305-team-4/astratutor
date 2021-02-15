package signalling

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Classes ...
type Classrooms struct {
	classes map[string]*Class
	mu      sync.Mutex
}

// Classroom ...
type Class struct {
	Code    string
	Members []*websocket.Conn
	mu      sync.Mutex
}

var classrooms = Classrooms{
	classes: map[string]*Class{},
}
