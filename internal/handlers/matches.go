package handlers

import (
	db "daterrr/internal/db/sqlc"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchHandler struct {
	store *db.SQLStore
}

func NewMatchHandler(store *db.SQLStore) *MatchHandler {
	return &MatchHandler{store: store}
}

// Lists all matches a user has gotten
func (m *MatchHandler) ListMatches(c *gin.Context) {
	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User id not found in context"})
	}

	// Assert to string first
	userIDStr, ok := userID.(string)
	if !ok {
		log.Printf("userID is not a string")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a string"})
	}

	// Then parse to UUID
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("invalid user ID format: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	matchList, err := m.store.GetUserMatches(c, UUIDToPgType(userUUID))
	if err != nil {
		log.Printf("ERROR: Something went wrong while fetching your matches %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while fetching your matches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matchList,
	})
}
