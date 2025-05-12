package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

// User represents a connected user
 type User struct {
	UserID   string          
	Conn *websocket.Conn 
}


type Lobby struct {
	ActiveUsers []*User   // Users waiting for a match
	Rooms        map[string]*ChatRoom // Active chat rooms
	Mutex        sync.Mutex // Synchronization for concurrent access
}

type ChatRoom struct {
	ID      string         // Unique identifier for the room
	Users   map[string]*User // Connected users
	Mutex   sync.Mutex     // Synchronization for concurrent access
	LogFile string         // Log file for messages
}