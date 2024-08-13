package main

import (
	"main-module/controllers"
	"main-module/initializers"
	"main-module/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	r.POST("/signup",controllers.Signup)
	r.POST("/login",controllers.Login)
	r.GET("/validate",middleware.RequireAuth,controllers.Validate)
	r.Run()
}
