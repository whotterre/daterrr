package handlers

import (
	"log"
	"net/http"
	"os"

	db "daterrr/internal/db/sqlc"
	"daterrr/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	store  *db.SQLStore
	logger *log.Logger
}

func NewChatHandler(store *db.SQLStore) *ChatHandler {
	// Initialize file-based logger
	logFile, err := os.OpenFile("chat_handler.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger := log.New(logFile, "CHAT_HANDLER: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &ChatHandler{store: store, logger: logger}
}

type CreateMessageRequest struct {
	ChatID     string `json:"chat_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	ReceiverID string `json:"receiver_id" binding:"required"`
}

func (h *ChatHandler) CreateMessage(c *gin.Context) {
	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.MustGet("userID").(string)
	senderUUID, _ := uuid.Parse(userID)
	chatUUID, _ := uuid.Parse(req.ChatID)

	message, err := h.store.CreateMessage(c, db.CreateMessageParams{
		ChatID:   UUIDToPgType(chatUUID),
		SenderID: UUIDToPgType(senderUUID),
		Content:  req.Content,
	})
	if err != nil {
		h.logger.Printf("Error creating message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Define the Message struct
	type Message struct {
		Type      string `json:"type"`
		ChatID    string `json:"chat_id"`
		Sender    string `json:"sender"`
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp"`
	}

	chatMessage := Message{
		Type:      "message",
		ChatID:    req.ChatID,
		Sender:    userID,
		Content:   req.Content,
		Timestamp: message.CreatedAt.Time.Unix(),
	}

	h.logger.Printf("Message created successfully: %+v", chatMessage)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent",
		"data":    chatMessage,
	})
}

func (h *ChatHandler) GetChatMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		h.logger.Printf("Invalid chat ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}
	h.logger.Print("UUID cast", chatUUID)
	h.logger.Print("PgType UUID", utils.UUIDToPgType(chatUUID))
	messages, err := h.store.GetChatMessages(c, utils.UUIDToPgType(chatUUID))
	if err != nil {
		h.logger.Printf("Failed to get messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	h.logger.Printf("Retrieved messages for chat ID %s", chatID)

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func (h *ChatHandler) GetUserChatsWithChatID(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Printf("Invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	chats, err := h.store.GetUserChatsWithChatID(c, utils.UUIDToPgType(userUUID))
	if err != nil {
		h.logger.Printf("Failed to get chats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	h.logger.Printf("Retrieved chats for user ID %s", userID)

	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
	})
}

func (h *ChatHandler) GetConversations(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.logger.Printf("Invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	conversations, err := h.store.GetConversationsForUser(c, utils.UUIDToPgType(userUUID))
	if err != nil {
		h.logger.Printf("Failed to get conversations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversations"})
		return
	}

	h.logger.Printf("Retrieved conversations for user ID %s", userID)

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
	})
}
