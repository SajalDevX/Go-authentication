package middleware

import (
	"fmt"
	"main-module/initializers"
	"main-module/models"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	// Get the cookie from the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret for validation
		return []byte(os.Getenv("SECRET")), nil
	})

	// Check if the token is valid and the claims are of type jwt.MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user based on the `sub` claim
		var user models.User
		if sub, ok := claims["sub"].(float64); ok {
			initializers.DB.First(&user, uint(sub))
			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach the user to the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
