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
	r.GET("/admin",middleware.RequireAuth, middleware.RoleMiddleware("admin"), controllers.AdminRoute)
	r.GET("/seller", middleware.RequireAuth,middleware.RoleMiddleware("seller"), controllers.SellerRoute)
	r.Run()
}
