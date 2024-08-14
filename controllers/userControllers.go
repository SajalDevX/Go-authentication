package controllers

import (
	"main-module/initializers"
	"main-module/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get the email/pass off the req body
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash Password"})
		return
	}

	// Create the user with a role
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
		Role:     body.Role, 
	}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}


func Login(c *gin.Context) {
	// Get the email and password from the request body
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"` // Fixed the binding tag from "email" to "password"
	}

	if err := c.BindJSON(&body); err != nil { // Changed c.Bind to c.BindJSON for better error handling
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Look up the user by email
	var user models.User
	result := initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the provided password with the saved hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	secret := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	// Send it back in a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true) // Path added for proper scope
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func AdminRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin!"})
}

func SellerRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome Seller!"})
}
