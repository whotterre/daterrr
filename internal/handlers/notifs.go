package handlers

import (
	"context"
	db "daterrr/internal/db/sqlc"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationHandler struct {
	store       *db.SQLStore
	connections map[string]*websocket.Conn
}

func NewNotificationHandler(store *db.SQLStore) *NotificationHandler {
	return &NotificationHandler{
		store:       store,
		connections: make(map[string]*websocket.Conn),
	}
}

type MatchNotification struct {
	Type      string      `json:"type"`
	MatchID   string      `json:"matchId"`
	UserID    pgtype.UUID `json:"userId"`
	FirstName string      `json:"firstName"`
	Age       int32       `json:"age"`
	ImageURL  string      `json:"imageUrl"`
	Timestamp int64       `json:"timestamp"`
}

func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	h.connections[userID] = conn
	defer delete(h.connections, userID)

	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			break
		}
	}
}

func (h *NotificationHandler) SendMatchNotification(user1ID, user2ID uuid.UUID) {
	ctx := context.Background()

	// Get user profiles
	user1, err := h.store.GetUserProfile(ctx, UUIDToPgType(user1ID))
	if err != nil {
		log.Printf("Error getting user1 profile: %v", err)
		return
	}

	user2, err := h.store.GetUserProfile(ctx, UUIDToPgType(user2ID))
	if err != nil {
		log.Printf("Error getting user2 profile: %v", err)
		return
	}

	// Prepare notifications
	notif1 := MatchNotification{
		Type:      "match",
		MatchID:   fmt.Sprintf("%s-%s", user1ID, user2ID),
		UserID:    UUIDToPgType(user2ID),
		FirstName: user2.FirstName,
		Age:       user2.Age,
		ImageURL:  user2.ImageUrl.String,
		Timestamp: time.Now().Unix(),
	}

	notif2 := MatchNotification{
		Type:      "match",
		MatchID:   fmt.Sprintf("%s-%s", user1ID, user2ID),
		UserID:    UUIDToPgType(user1ID),
		FirstName: user1.FirstName,
		Age:       user1.Age,
		ImageURL:  user1.ImageUrl.String,
		Timestamp: time.Now().Unix(),
	}

	// Send to user1 if connected
	if conn, ok := h.connections[user1ID.String()]; ok {
		if err := conn.WriteJSON(notif1); err != nil {
			log.Printf("Error sending WS to user1: %v", err)
		}
	}

	// Send to user2 if connected
	if conn, ok := h.connections[user2ID.String()]; ok {
		if err := conn.WriteJSON(notif2); err != nil {
			log.Printf("Error sending WS to user2: %v", err)
		}
	}

	// Persist notifications
	h.storeNotification(ctx, user1ID.String(), notif1)
	h.storeNotification(ctx, user2ID.String(), notif2)
}

func (h *NotificationHandler) storeNotification(ctx context.Context, userID string, notif MatchNotification) {
	notifData, err := json.Marshal(notif)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return
	}

	uuid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Something went wrong in converting the userID from string to uuid")

	}
	_, err = h.store.CreateNotification(ctx, db.CreateNotificationParams{
		UserID:   	UUIDToPgType(uuid),
		Type:      "match",
		Data:      notifData,
	})
	if err != nil {
		log.Printf("Error saving notification: %v", err)
	}
}

