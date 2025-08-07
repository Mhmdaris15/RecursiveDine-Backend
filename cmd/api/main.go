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
// @host localhost:8002
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
	userService := services.NewUserService(userRepo)
	tableService := services.NewTableService(tableRepo)
	menuService := services.NewMenuService(menuRepo)
	orderService := services.NewOrderService(orderRepo, menuRepo)
	paymentService := services.NewPaymentService(paymentRepo, orderRepo, cfg)
	kitchenService := services.NewKitchenService(orderRepo)
	seedService := services.NewSeedService(db, userRepo, tableRepo, menuRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	tableController := controllers.NewTableController(tableService)
	menuController := controllers.NewMenuController(menuService)
	orderController := controllers.NewOrderController(orderService, authService)
	paymentController := controllers.NewPaymentController(paymentService)
	kitchenController := controllers.NewKitchenController(kitchenService)
	seedController := controllers.NewSeedController(seedService)
	
	// Initialize CRUD controllers
	userController := controllers.NewUserController(userService)
	orderManagementController := controllers.NewOrderManagementController(orderService)
	paymentManagementController := controllers.NewPaymentManagementController(paymentService)

	// Setup router
	router := setupRouter(cfg, authController, tableController, menuController, orderController, paymentController, kitchenController, userController, orderManagementController, paymentManagementController, seedController)

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
	fmt.Printf("Debug - Database Config Values:\n")
	fmt.Printf("Host: '%s'\n", cfg.DBHost)
	fmt.Printf("User: '%s'\n", cfg.DBUser)
	fmt.Printf("Password: '%s'\n", cfg.DBPassword)
	fmt.Printf("Database: '%s'\n", cfg.DBName)
	fmt.Printf("Port: '%s'\n", cfg.DBPort)
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	
	fmt.Printf("DSN: %s\n", dsn)

	// Retry database connection up to 10 times with increasing delays
	var db *gorm.DB
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Test the connection
			sqlDB, sqlErr := db.DB()
			if sqlErr == nil {
				pingErr := sqlDB.Ping()
				if pingErr == nil {
					fmt.Printf("Successfully connected to database on attempt %d\n", i+1)
					break
				}
				err = pingErr
			} else {
				err = sqlErr
			}
		}
		
		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * 2 * time.Second
			fmt.Printf("Database connection attempt %d failed: %v. Retrying in %v...\n", i+1, err, waitTime)
			time.Sleep(waitTime)
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// // Auto-migrate models
	// err = db.AutoMigrate(
	// 	&repositories.User{},
	// 	&repositories.Table{},
	// 	&repositories.MenuCategory{},
	// 	&repositories.MenuItem{},
	// 	&repositories.Order{},
	// 	&repositories.OrderItem{},
	// 	&repositories.Payment{},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return db, nil
}

func setupRouter(cfg *config.Config, authController *controllers.AuthController, tableController *controllers.TableController, menuController *controllers.MenuController, orderController *controllers.OrderController, paymentController *controllers.PaymentController, kitchenController *controllers.KitchenController, userController *controllers.UserController, orderManagementController *controllers.OrderManagementController, paymentManagementController *controllers.PaymentManagementController, seedController *controllers.SeedController) *gin.Engine {
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
			auth.POST("/logout", middleware.AuthMiddleware(cfg), authController.Logout)
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
		orders.Use(middleware.AuthMiddleware(cfg))
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("/:id", orderController.GetOrder)
			orders.GET("", orderController.GetOrders)
			orders.PATCH("/:id/status", middleware.RoleMiddleware("staff", "admin"), orderController.UpdateOrderStatus)
		}

		// Payment routes
		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware(cfg))
		{
			payments.POST("/qris", paymentController.InitiateQRISPayment)
			payments.POST("/verify", paymentController.VerifyPayment)
			payments.GET("/status/:payment_id", paymentController.GetPaymentStatus)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// User management
			users := admin.Group("/users")
			{
				users.GET("", userController.GetAllUsers)
				users.GET("/:id", userController.GetUserByID)
				users.POST("", userController.CreateUser)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)
			}

			// Table management
			tables := admin.Group("/tables")
			{
				tables.GET("", tableController.GetAllTables)
				tables.POST("", tableController.CreateTable)
				tables.PUT("/:id", tableController.UpdateTable)
				tables.DELETE("/:id", tableController.DeleteTable)
				tables.PATCH("/:id/status", tableController.UpdateTableAvailability)
			}

			// Menu management
			menuAdmin := admin.Group("/menu")
			{
				// Category management
				menuAdmin.POST("/categories", menuController.CreateCategory)
				menuAdmin.PUT("/categories/:id", menuController.UpdateCategory)
				menuAdmin.DELETE("/categories/:id", menuController.DeleteCategory)

				// Menu item management
				menuAdmin.POST("/items", menuController.CreateMenuItem)
				menuAdmin.PUT("/items/:id", menuController.UpdateMenuItem)
				menuAdmin.DELETE("/items/:id", menuController.DeleteMenuItem)
				menuAdmin.PATCH("/items/:id/availability", menuController.UpdateMenuItemAvailability)
			}

			// Order management
			orders := admin.Group("/orders")
			{
				orders.GET("", orderManagementController.GetAllOrders)
				orders.GET("/:id", orderManagementController.GetOrderByID)
				orders.PATCH("/:id/status", orderManagementController.UpdateOrderStatus)
				orders.PUT("/:id/items", orderManagementController.UpdateOrderItems)
				orders.DELETE("/:id", orderManagementController.DeleteOrder)
				orders.GET("/statistics", orderManagementController.GetOrderStatistics)
				orders.GET("/revenue", orderManagementController.GetDailyRevenue)
			}

			// Payment management
			paymentAdmin := admin.Group("/payments")
			{
				paymentAdmin.GET("", paymentManagementController.GetAllPayments)
				paymentAdmin.GET("/:id", paymentManagementController.GetPaymentByID)
				paymentAdmin.PATCH("/:id/status", paymentManagementController.UpdatePaymentStatus)
				paymentAdmin.POST("/:id/refund", paymentManagementController.ProcessRefund)
				paymentAdmin.DELETE("/:id", paymentManagementController.DeletePayment)
				paymentAdmin.GET("/statistics", paymentManagementController.GetPaymentStatistics)
				paymentAdmin.GET("/revenue", paymentManagementController.GetDailyRevenueByPayment)
			}

			// Database seeding routes
			seed := admin.Group("/seed")
			{
				seed.POST("", seedController.SeedDatabase)
				seed.DELETE("", seedController.ClearDatabase)
				seed.GET("/status", seedController.GetSeedStatus)
			}
		}

		// Staff routes (staff and admin access)
		staff := api.Group("/staff")
		staff.Use(middleware.AuthMiddleware(cfg))
		staff.Use(middleware.RoleMiddleware("staff", "admin"))
		{
			// Order management for staff
			orders := staff.Group("/orders")
			{
				orders.GET("", orderManagementController.GetAllOrders)
				orders.GET("/:id", orderManagementController.GetOrderByID)
				orders.PATCH("/:id/status", orderManagementController.UpdateOrderStatus)
			}

			// Menu availability updates
			menu := staff.Group("/menu")
			{
				menu.PATCH("/items/:id/availability", menuController.UpdateMenuItemAvailability)
			}
		}

		// Cashier routes (cashier and admin access)
		cashier := api.Group("/cashier")
		cashier.Use(middleware.AuthMiddleware(cfg))
		cashier.Use(middleware.RoleMiddleware("cashier", "admin"))
		{
			// Cash payment processing
			payments := cashier.Group("/payments")
			{
				payments.POST("/cash", paymentController.ProcessCashPayment)
				payments.GET("", paymentManagementController.GetAllPayments)
				payments.GET("/:id", paymentManagementController.GetPaymentByID)
				payments.POST("/:id/refund", paymentManagementController.ProcessRefund)
				payments.POST("/reconcile", paymentManagementController.ReconcileCashPayments)
				payments.GET("/statistics", paymentManagementController.GetPaymentStatistics)
			}
		}
	}

	// WebSocket for kitchen updates
	router.GET("/kitchen/updates", middleware.WSAuthMiddleware(cfg), kitchenController.HandleWebSocket)

	return router
}
