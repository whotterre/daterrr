package handlers

import (
	"context"
	db "daterrr/internal/db/sqlc"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)
type SwipeHandler struct {
	store               *db.SQLStore
	notificationHandler *NotificationHandler
}

func NewSwipeHandler(store *db.SQLStore, notificationHandler *NotificationHandler) *SwipeHandler {
	return &SwipeHandler{
		store:               store,
		notificationHandler: notificationHandler,
	}
}

type SwipeRequest struct {
	SwipeeID string `json:"swipeeID"`
}

func UUIDToPgType(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func (s *SwipeHandler) CreateFeed(c *gin.Context) {
    // Get userID from context
    userID, exists := c.Get("userID")
    if !exists {
        log.Printf("ERROR: userID not found in context")
        c.JSON(http.StatusBadRequest, gin.H{"error": "User authentication required"})
        return
    }

    // Type assertion
    userIDStr, ok := userID.(string)
    if !ok {
        log.Printf("ERROR: userID is not a string (got %T)", userID)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user format"})
        return
    }

    // Convert to UUID
    userUUID, err := uuid.Parse(userIDStr)
    if err != nil {
        log.Printf("ERROR: Failed to parse UUID: %v (input: %s)", err, userIDStr)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID format",
            "details": err.Error(),
        })
        return
    }

    // Execute query with context timeout
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5* time.Second)
    defer cancel()

    feedItems, err := s.store.GenerateFeed(ctx, UUIDToPgType(userUUID))
    if err != nil {
        log.Printf("ERROR: GenerateFeed failed: %v", err)
        
        if errors.Is(err, context.DeadlineExceeded) {
            c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to generate feed",
                "details": err.Error(),
            })
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "profilePool": feedItems,
    })
}

func (s *SwipeHandler) HandleSwipe(c *gin.Context) {
    // Initialize logging properly
    var req SwipeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Printf("Invalid request body: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    
    // Get current user ID from context
    swiperID, exists := c.Get("userID")
    if !exists {
        log.Println("userID not found in context")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
        return
    }
    
    // Convert IDs to UUID type
    swiperUUID, err := uuid.Parse(swiperID.(string))
    if err != nil {
        log.Printf("Invalid swiper UUID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    swipeeUUID, err := uuid.Parse(req.SwipeeID)
    if err != nil {
        log.Printf("Invalid swipee UUID: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swipee ID"})
        return
    }
    
    var matchID interface{}
    var isMutual bool
    
    // Execute all DB operations in a transaction
    err = s.store.ExecTx(c, func(q db.Querier) error {
        // Check if swipee exists
        if _, err := q.GetUserByID(c, UUIDToPgType(swipeeUUID)); err != nil {
            if err == pgx.ErrNoRows {
                log.Printf("Swipee not found: %v", swipeeUUID)
                return fmt.Errorf("user profile not found")
            }
            log.Printf("Error checking user profile: %v", err)
            return fmt.Errorf("failed to check user profile: %w", err)
        }
        
        // Record the swipe
        if err := q.NewSwipe(c, db.NewSwipeParams{
            SwiperID: UUIDToPgType(swiperUUID),
            SwipeeID: UUIDToPgType(swipeeUUID),
        }); err != nil {
            log.Printf("Failed to record swipe: %v", err)
            log.Print(swipeeUUID, swiperUUID)
            return fmt.Errorf("failed to record swipe: %w", err)
        }
        
        // Check for mutual swipes
        mutual, err := q.CheckMutualSwipe(c, db.CheckMutualSwipeParams{
            SwiperID: UUIDToPgType(swiperUUID),  
            SwipeeID: UUIDToPgType(swipeeUUID),
        })
        if err != nil {
            log.Printf("Failed to check mutual swipe: %v", err)
            return fmt.Errorf("failed to check mutual swipe: %w", err)
        }
        
        isMutual = mutual
        
        if mutual {
            match, err := q.CreateMatch(c, db.CreateMatchParams{
                Column1: UUIDToPgType(swiperUUID),
                Column2: UUIDToPgType(swipeeUUID),
            })
            if err != nil {
                // Handle "no rows" error - can occur if match already exists
                if errors.Is(err, pgx.ErrNoRows) {
                    log.Printf("Match may already exist, continuing: %v", err)
                    // Try to find existing match ID
                    existingID, findErr := q.FindExistingMatch(c, db.FindExistingMatchParams{
                        User1ID: UUIDToPgType(minUUID(swiperUUID, swipeeUUID)),
                        User2ID: UUIDToPgType(maxUUID(swiperUUID, swipeeUUID)),
                    })
                    if findErr == nil {
                        matchID = existingID
                        return nil
                    }
                    return fmt.Errorf("failed to create or find match: %w", err)
                }
                
                log.Printf("Failed to create match: %v", err)
                return fmt.Errorf("failed to create match: %w", err)
            }
            matchID = match
        }
        
        return nil
    })
    
    if err != nil {
        // Handle different types of errors
        if err.Error() == "user profile not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
        return
    }
    
    // After successful transaction, send notifications if there was a match
    if isMutual {
        // Send real-time notifications
        s.notificationHandler.SendMatchNotification(swiperUUID, swipeeUUID)
        
        c.JSON(http.StatusOK, gin.H{
            "message": "It's a match!",
            "matchID": matchID,
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Swipe recorded"})
}

// Helper functions to determine min and max UUIDs
func minUUID(a, b uuid.UUID) uuid.UUID {
    if a.String() < b.String() {
        return a
    }
    return b
}

func maxUUID(a, b uuid.UUID) uuid.UUID {
    if a.String() > b.String() {
        return a
    }
    return b
}