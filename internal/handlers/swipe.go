package handlers

import (
	db "daterrr/internal/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type SwipeHandler struct {
	store              *db.SQLStore
	notificationHandler *NotificationHandler
}

func NewSwipeHandler(store *db.SQLStore, notificationHandler *NotificationHandler) *SwipeHandler {
	return &SwipeHandler{
		store:              store,
		notificationHandler: notificationHandler,
	}
}

type SwipeRequest struct {
	SwipeeID string `json:"swipeeID"` // The ID of the user being swiped on
}

func UUIDToPgType(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func (s *SwipeHandler) HandleSwipe(c *gin.Context) {
	var req SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get current user ID from context
	swiperID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Convert IDs to UUID type
	swiperUUID, err := uuid.Parse(swiperID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	swipeeUUID, err := uuid.Parse(req.SwipeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swipee ID"})
		return
	}
	
	// Check if swipee exists
	_, err = s.store.GetUserByID(c, UUIDToPgType(swipeeUUID))
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user profile"})
		return
	}

	// Record the swipe (assuming all swipes are right swipes/likes)
	err = s.store.NewSwipe(c, db.NewSwipeParams{
		SwiperID: UUIDToPgType(swiperUUID),
		SwipeeID: UUIDToPgType(swipeeUUID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record swipe"})
		return
	}

	// Check for mutual swipe
	mutual, err := s.store.CheckMutualSwipe(c, db.CheckMutualSwipeParams{
		SwiperID: UUIDToPgType(swiperUUID),
		SwipeeID: UUIDToPgType(swipeeUUID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check mutual swipe"})
		return
	}

	// If mutual swipe exists, create a match
	if mutual {
		match, err := s.store.CreateMatch(c, db.CreateMatchParams{
			Column1: UUIDToPgType(swiperUUID),
			Column2: UUIDToPgType(swipeeUUID),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create match"})
			return
		}

		// Send real-time notifications to both users
		s.notificationHandler.SendMatchNotification(swiperUUID, swipeeUUID)

		c.JSON(http.StatusOK, gin.H{
			"message": "It's a match!",
			"matchID": match,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Swipe recorded"})
}