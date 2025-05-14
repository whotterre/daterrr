package handlers

import (
	db "daterrr/internal/db/sqlc"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	store *db.SQLStore
}


func NewProfileHandler(store *db.SQLStore) *ProfileHandler {
	return &ProfileHandler{store: store}
}

func (p *ProfileHandler) GetUserProfile(c *gin.Context){
	// Check store	
	userID, exists := c.Get("userID")
    if !exists {
        log.Printf("ERROR: userID not found in context")
        c.JSON(http.StatusBadRequest, gin.H{"error": "User authentication required"})
        return
    }

	userIDStr, ok := userID.(string)
    if !ok {
        log.Printf("ERROR: userID is not a string (got %T)", userID)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
        return
    }

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Print("Failed to convert string to uuid because", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on our end"})
        return
	}

	profile, err := p.store.GetUserProfile(c, UUIDToPgType(userUUID))
	if err != nil {
		log.Printf("Couldn't fetch user profile because %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong our end while fetching your user profile",
		})
		return
	}
	

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
	})

}