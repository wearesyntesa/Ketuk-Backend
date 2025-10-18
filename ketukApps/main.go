package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ketukApps/config"
	"ketukApps/internal/database"
	"ketukApps/internal/handlers"
	"ketukApps/internal/middleware"
	"ketukApps/internal/queue"
	"ketukApps/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	db := database.GetDB()

	// Initialize RabbitMQ
	if err := queue.NewRabbitMQConnection(cfg); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer queue.CloseRabbitMQ()

	// Initialize services
	userService := services.NewUserService(db)
	ticketService := services.NewTicketService(db)
	scheduleService := services.NewScheduleService(db)

	// Start the worker with ticket service and schedule service
	go func() {
		if err := queue.SchduleWorker("schedule", ticketService, scheduleService); err != nil {
			log.Fatalf("Failed to start schedule worker: %v", err)
		}
	}()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	tickets := handlers.NewTicketHandler(ticketService)

	// Setup Gin router
	router := setupRouter(userHandler, tickets)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Printf("🚀 Server starting on http://%s", address)
	log.Printf("📚 API Documentation: http://%s/api/docs", address)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(userHandler *handlers.UserHandler, ticketHandler *handlers.TicketHandler) *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode) // Change to gin.DebugMode for development

	// Create router with default middleware
	router := gin.New()

	// Add custom middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// API routes group
	api := router.Group("/api")
	{
		// Users endpoints
		users := api.Group("/users")
		{
			users.GET("/v1", userHandler.GetAllUsers)
			users.GET("/v1/:id", userHandler.GetUserByID)
			users.POST("v1", userHandler.CreateUser)
			users.PUT("/v1/:id", userHandler.UpdateUser)
			users.DELETE("/v1/:id", userHandler.DeleteUser)
		}

		tickets := api.Group("/tickets")
		{
			// Ticket endpoints would go here
			tickets.GET("/v1", ticketHandler.GetAllTickets)
			tickets.GET("/v1/:id", ticketHandler.GetTicketByID)
			tickets.POST("/v1", ticketHandler.CreateTicket)
			tickets.PUT("/v1/:id", ticketHandler.UpdateTicket)
			tickets.DELETE("/v1/:id", ticketHandler.DeleteTicket)
			tickets.PATCH("/v1/:id/status", ticketHandler.UpdateTicketStatus)
			tickets.POST("/v1/bulk-status", ticketHandler.BulkUpdateStatus)
		}
	}

	// Serve API documentation (simple endpoint)
	router.GET("/api/docs", func(c *gin.Context) {
		docs := `
		<h1>KetukApps API Documentation</h1>
		<h2>Available Endpoints:</h2>
		<ul>
			<li><strong>GET /health</strong> - Health check</li>
			<li><strong>GET /api/users</strong> - Get all users</li>
			<li><strong>GET /api/users/:id</strong> - Get user by ID</li>
			<li><strong>POST /api/users</strong> - Create new user</li>
			<li><strong>PUT /api/users/:id</strong> - Update user</li>
			<li><strong>DELETE /api/users/:id</strong> - Delete user</li>
		</ul>
		<h3>Example User JSON:</h3>
		<pre>
{
  "name": "John Doe",
  "email": "john@example.com"
}
		</pre>
		`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(docs))
	})

	return router
}
