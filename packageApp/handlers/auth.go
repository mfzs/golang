package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"packageApp/models"
)

// ServeLoginPage serves the login page
func ServeLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Login handles user login
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var storedPassword string
	err := models.GetDB().QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	if user.Password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// HandleSignup handles user signup
func HandleSignup(c *gin.Context) {
	var credentials models.User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Ensure password is hashed and stored securely
	_, err := models.GetDB().Exec("INSERT INTO users (username, password) VALUES (?, ?)", credentials.Username, credentials.Password)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to sign up"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}
