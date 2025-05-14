package handlers

import (
	db "daterrr/internal/db/sqlc"
	"daterrr/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type IncomingMessage struct {
	Type      string `json:"type"`
	ChatID    string `json:"chat_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

var lobby = &models.Lobby{
	ActiveUsers: make([]*models.User, 0),
	Rooms:       make(map[string]*models.ChatRoom),
}

// HandleWebSocket manages user connections and messaging between matches
func HandleWebSocket(c *gin.Context, store *db.SQLStore) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("[ERROR] Failed to upgrade WebSocket connection:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	fmt.Println("[CONNECTED] WebSocket connection established")

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	userID := c.Query("userId")
	if userID == "" {
		fmt.Println("[ERROR] Missing userId query parameter")
		http.Error(c.Writer, "User ID is required", http.StatusBadRequest)
		return
	}

	matchID := c.Query("matchId")
	if matchID == "" {
		fmt.Println("[ERROR] Missing matchId query parameter")
		http.Error(c.Writer, "Match ID is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("[INFO] User %s connected to match %s\n", userID, matchID)

	user := &models.User{
		UserID: userID,
		Conn:   conn,
	}

	lobby.Mutex.Lock()
	defer lobby.Mutex.Unlock()

	// Add user to active users
	lobby.ActiveUsers = append(lobby.ActiveUsers, user)
	// Notify the user they are connected
	room, exists := lobby.Rooms[matchID]
	if !exists {
		room = &models.ChatRoom{
			ID:    matchID,
			Users: make(map[string]*models.User),
		}
		lobby.Rooms[matchID] = room
	}

	// Add user to the room
	room.Users[userID] = user

	// Notify the user they've joined the chat
	welcomeMsg := map[string]string{
		"sender":  "Server",
		"message": fmt.Sprintf("You are now connected to your match chat"),
		"type": "message",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}
	jsonMsg, _ := json.Marshal(welcomeMsg)
	conn.WriteMessage(websocket.TextMessage, jsonMsg)

	// Notify the other user in the room (if present)
	if len(room.Users) == 2 {
		for _, u := range room.Users {
			if u.UserID != userID {
				notification := map[string]string{
					"sender":  "Server",
					"message": "Your match has joined the chat",
					"type":    "message",
					"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
				}
				jsonNotif, _ := json.Marshal(notification)
				u.Conn.WriteMessage(websocket.TextMessage, jsonNotif)
			}
		}
	}

	// Start listening for messages from this user
	go listenForMessages(c, user, room, store)
}

// Listens for messages from a user and broadcasts them to their match
func listenForMessages(context *gin.Context, user *models.User, room *models.ChatRoom, store *db.SQLStore) {
	defer func() {
		removeUserFromRoom(user.UserID)
		user.Conn.Close()
	}()

	for {
		_, msg, err := user.Conn.ReadMessage()
		if err != nil {
			fmt.Println("[DISCONNECT] User disconnected:", user.UserID, "Error:", err)
			break
		}

		var messageData IncomingMessage
		if err := json.Unmarshal(msg, &messageData); err != nil {
			fmt.Println("[ERROR] Failed to parse message:", err)
			continue
		}

		if messageData.Type == "leave" {
			fmt.Println("[USER LEFT] User intentionally left:", user.UserID)
			break
		}

		fmt.Println("[MESSAGE RECEIVED] From:", user.UserID, "Message:", string(msg))
		broadcastMessage(context, room, user.UserID, msg, store)
	}
}

// Broadcasts messages within a room (to the matched user)
func broadcastMessage(c *gin.Context, room *models.ChatRoom, senderID string, msg []byte, store *db.SQLStore) {
	lobby.Mutex.Lock()
	defer lobby.Mutex.Unlock()

	var parsedMsg IncomingMessage
	if err := json.Unmarshal(msg, &parsedMsg); err != nil {
		fmt.Println("[ERROR] Failed to parse incoming message:", err)
		return
	}
	if parsedMsg.Type == "typing" {
		typingNotification := map[string]string{
			"type":    "typing",
			"sender":  senderID,
			"message": fmt.Sprintf("User %s is typing...", senderID),
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		}
		jsonTyping, _ := json.Marshal(typingNotification)

		for _, user := range room.Users {
			if user.UserID != senderID {
				err := user.Conn.WriteMessage(websocket.TextMessage, jsonTyping)
				if err != nil {
					fmt.Println("[ERROR] Failed to send typing notification to user:", user.UserID, err)
					removeUserFromRoom(user.UserID)
				}
			}
		}
		fmt.Printf("[TYPING] User %s is typing in room %s\n", senderID, room.ID)
		return
	}
	messageData := map[string]string{
		"sender":  senderID,
		"message": parsedMsg.Content,
		"type":    "message",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}
	/* Persist the messages */

	// Convert senderID to UUID
	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		fmt.Println("[ERROR] Invalid sender ID:", senderID, err)
		return
	}
	// Convert room ID to UUID
	fmt.Println("[INFO] Room ID:", room.ID)
	log.Print(room.ID)
	roomUUID, err := uuid.Parse(room.ID)
	if err != nil {
		fmt.Println("[ERROR] Invalid room ID:", room.ID, err)
		return
	}
	// Store the message in the database
	message, err := store.CreateMessage(c, db.CreateMessageParams{
		ChatID:   UUIDToPgType(roomUUID),
		SenderID: UUIDToPgType(senderUUID),
		Content:  parsedMsg.Content,
	})
	if err != nil {
		fmt.Println("[ERROR] Failed to store message in database:", message, err)
		return
	}

	jsonMsg, _ := json.Marshal(messageData)

	for _, user := range room.Users {
		if user.UserID != senderID {
			err := user.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
			if err != nil {
				fmt.Println("[ERROR] Failed to send message to user:", user.UserID, err)
				removeUserFromRoom(user.UserID)
			}
		}
	}
}

// Removes a user when they disconnect
func removeUserFromRoom(userID string) {
	lobby.Mutex.Lock()
	defer lobby.Mutex.Unlock()

	// Remove from active users
	for i, user := range lobby.ActiveUsers {
		if user.UserID == userID {
			lobby.ActiveUsers = append(lobby.ActiveUsers[:i], lobby.ActiveUsers[i+1:]...)
			break
		}
	}

	// Find and remove from rooms
	for roomID, room := range lobby.Rooms {
		if _, exists := room.Users[userID]; exists {
			delete(room.Users, userID)
			fmt.Println("[USER LEFT] User", userID, "left room", roomID)

			// Notify the other user in the room
			for _, remainingUser := range room.Users {
				disconnectMsg := map[string]string{
					"sender":  "Server",
					"message": "Your match has disconnected from the chat.",
					"type":    "system",
					"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
				}
				jsonMsg, _ := json.Marshal(disconnectMsg)
				remainingUser.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
			}

			// Remove empty rooms
			if len(room.Users) == 0 {
				delete(lobby.Rooms, roomID)
				fmt.Println("[ROOM DELETED] Empty room", roomID, "removed")
			}
			break
		}
	}
}
