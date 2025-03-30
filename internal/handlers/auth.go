package handlers

import (
	"daterrr/internal/db/sqlc"
	"daterrr/internal/utils"
	"daterrr/pkg/auth"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthHandler struct {
	store *db.SQLStore
}

func NewAuthHandler(store *db.SQLStore) *AuthHandler {
	return &AuthHandler{store: store}
}

type CreateNewUserParams struct {
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Bio       pgtype.Text `json:"bio"`
	Gender    string      `json:"gender"`
	Age       int32       `json:"age"`
	ImageUrl  pgtype.Text `json:"image_url"`
	Interests []string    `json:"interests"`
}

// Registers a new user
func (a *AuthHandler) RegisterUser(c *gin.Context) {
	// Bind request body
	var req CreateNewUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// _, err := a.store.GetUserByEmail(c, req.Email)
	// switch {
	// case errors.Is(err, pgx.ErrNoRows):
	// 	c.JSON(http.StatusConflict, gin.H{
	// 		"error": "Email already registered. Please log in.",
	// 	})
	// 	return
	// case err != nil && !errors.Is(err, pgx.ErrNoRows):
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Could not verify email availability",
	// 	})
	// 	return
	// }
	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Couldn't hash password",
		})
		return
	}

	// Create user record
	newUser, err := a.store.CreateNewUser(c, db.CreateNewUserParams{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		Gender:    req.Gender,
		Age:       req.Age,
		ImageUrl:  req.ImageUrl,
		Interests: req.Interests,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong while creating new user record in the database",
			"error":   err.Error(),
		})
		return
	}
	// Send a welcome email
	configPath := filepath.Join("../../")
	config, err := utils.LoadConfig(configPath)

	emailSender := utils.NewGmailSender("Daterrr", config.EmailAddr, config.GmailKey)

	subject := "Welcome to Daterrr!"
	content := fmt.Sprintf(`
	<h1>Welcome to Daterrr application, </h1>, %s\n
	<p> Thanks for joining us. Have a nice time.</p>
	`, req.FirstName)
	to := []string{req.Email}

	emailSender.SendEmail(subject, content, to, nil, nil, nil)
	c.JSON(http.StatusOK, gin.H{
		"message": "User profile created successfully",
		"user":    newUser,
	})
}
