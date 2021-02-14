package signalling

import "github.com/gorilla/websocket"

// Classroom ...
type Classroom struct {
	Code    string
	Members []*websocket.Conn
}

var classes = map[string]*Classroom{}
