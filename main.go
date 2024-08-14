package main

import (
	"main-module/controllers"
	"main-module/initializers"
	"main-module/middleware"
	roles "main-module/middleware/role"

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
	r.GET("/admin",middleware.RequireAuth, middleware.RoleMiddleware(roles.Admin), controllers.AdminRoute)
	r.GET("/seller", middleware.RequireAuth,middleware.RoleMiddleware(roles.Seller), controllers.SellerRoute)
	r.Run()
}
