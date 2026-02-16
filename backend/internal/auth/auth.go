package auth

import (
	"capuchin/internal/database"
	"capuchin/internal/models"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(c *gin.Context) {

	var input models.User

	// Validate request body
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.PasswordHash),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "password hashing failed",
		})
		return
	}

	input.PasswordHash = string(hash)

	// Save user to DB
	result := database.DB.Create(&input)

	if result.Error != nil {

		if strings.Contains(result.Error.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "user already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user creation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
	})
}

func LoginHandler(c *gin.Context) {

	var input models.User
	var user models.User

	// Parse request body
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// Find user by email
	result := database.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}

	// Compare password hash
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.PasswordHash),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(user.UserID, user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "token generation failed",
		})
		return
	}

	// Return token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
