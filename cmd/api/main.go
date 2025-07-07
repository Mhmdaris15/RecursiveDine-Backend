package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"recursiveDine/internal/config"
	"recursiveDine/internal/controllers"
	"recursiveDine/internal/middleware"
	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title RecursiveDine API
// @version 1.0
// @description A secure restaurant ordering system API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	tableRepo := repositories.NewTableRepository(db)
	menuRepo := repositories.NewMenuRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	tableService := services.NewTableService(tableRepo)
	menuService := services.NewMenuService(menuRepo)
	orderService := services.NewOrderService(orderRepo, menuRepo)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo, cfg)
	kitchenService := services.NewKitchenService(orderRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	tableController := controllers.NewTableController(tableService)
	menuController := controllers.NewMenuController(menuService)
	orderController := controllers.NewOrderController(orderService, authService)
	paymentController := controllers.NewPaymentController(paymentService)
	kitchenController := controllers.NewKitchenController(kitchenService)

	// Setup router
	router := setupRouter(cfg, authController, tableController, menuController, orderController, paymentController, kitchenController)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&repositories.User{},
		&repositories.Table{},
		&repositories.MenuCategory{},
		&repositories.MenuItem{},
		&repositories.Order{},
		&repositories.OrderItem{},
		&repositories.Payment{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupRouter(cfg *config.Config, authController *controllers.AuthController, tableController *controllers.TableController, menuController *controllers.MenuController, orderController *controllers.OrderController, paymentController *controllers.PaymentController, kitchenController *controllers.KitchenController) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit())

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(), authController.Logout)
		}

		// Table routes
		tables := api.Group("/tables")
		{
			tables.GET("/:qr_code", tableController.GetTableByQRCode)
		}

		// Menu routes
		menu := api.Group("/menu")
		{
			menu.GET("", menuController.GetMenu)
			menu.GET("/categories", menuController.GetCategories)
			menu.GET("/items/search", menuController.SearchItems)
		}

		// Order routes
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("/:id", orderController.GetOrder)
			orders.GET("", orderController.GetOrders)
			orders.PATCH("/:id/status", middleware.RoleMiddleware("staff", "admin"), orderController.UpdateOrderStatus)
		}

		// Payment routes
		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware())
		{
			payments.POST("/qris", paymentController.InitiateQRISPayment)
			payments.POST("/verify", paymentController.VerifyPayment)
			payments.GET("/status/:payment_id", paymentController.GetPaymentStatus)
		}
	}

	// WebSocket for kitchen updates
	router.GET("/kitchen/updates", middleware.WSAuthMiddleware(), kitchenController.HandleWebSocket)

	return router
}
