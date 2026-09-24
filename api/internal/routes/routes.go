package routes

import (
	"atonixcorp/api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authController *controllers.AuthController,
	contactController *controllers.ContactController,
	blogController *controllers.BlogController,
) {

	// API Versioning
	api := r.Group("/api/v1")

	// -------------------------
	// AUTH ROUTES
	// -------------------------
	auth := api.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/login", authController.Login)
	}

	// -------------------------
	// CONTACT ROUTES
	// -------------------------
	contact := api.Group("/contact")
	{
		contact.POST("/submit", contactController.Submit)
		contact.GET("/messages", contactController.GetAll) // Admin only (JWT later)
	}

	// -------------------------
	// BLOG ROUTES
	// -------------------------
	blog := api.Group("/blog")
	{
		blog.POST("/", blogController.Create)
		blog.PUT("/", blogController.Update)
		blog.DELETE("/:id", blogController.Delete)
		blog.GET("/:id", blogController.GetByID)
		blog.GET("/", blogController.GetAll)
	}
}
