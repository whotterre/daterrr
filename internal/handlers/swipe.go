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
	store *db.SQLStore
}

func NewSwipeHandler(store *db.SQLStore) *SwipeHandler {
	return &SwipeHandler{store: store}
}

type SwipeRequest struct {
	SwipeeID string `json:"swipeeID"`
}

func UUIDToPgType(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// Swipe on frontend
// UserID in request
// Add to swiped table
func (s *SwipeHandler) HandleSwipe(c *gin.Context) {
	var req SwipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	swiperID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	// Convert to ids to UUID type
	swiperUUID, err := uuid.Parse(swiperID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	
	swipeeUUID, err := uuid.Parse(req.SwipeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid swipee ID",
		})
		return
	}
	
	// Check if swipee exists
	_, err = s.store.GetUserByID(c, UUIDToPgType(swipeeUUID))
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check user profile",
		})
		return
	}

	// Insert into the swipes table
	err = s.store.NewSwipe(c, db.NewSwipeParams{
		SwiperID: UUIDToPgType(swiperUUID),
		SwipeeID: UUIDToPgType(swipeeUUID),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to create match.",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Swipe successfully recorded",
	})
	
}
