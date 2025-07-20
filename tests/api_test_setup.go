package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"recursiveDine/internal/config"
	"recursiveDine/internal/controllers"
	"recursiveDine/internal/middleware"
	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type APITestSuite struct {
	suite.Suite
	router      *gin.Engine
	db          *gorm.DB
	adminToken  string
	staffToken  string
	cashierToken string
	customerToken string
}

// SetupSuite runs before all tests
func (suite *APITestSuite) SetupSuite() {
	fmt.Println("SetupSuite starting...")
	
	// Set test environment
	os.Setenv("ENVIRONMENT", "test")
	gin.SetMode(gin.TestMode)

	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		suite.Require().NoError(err)
	}
	fmt.Println("Configuration loaded successfully")

	// Use test database
	cfg.DBName = cfg.DBName + "_test"
	fmt.Printf("Using test database: %s\n", cfg.DBName)

	// Initialize test database
	suite.db, err = setupTestDatabase(cfg)
	if err != nil {
		fmt.Printf("Failed to setup database: %v\n", err)
		suite.Require().NoError(err)
	}
	fmt.Println("Database setup completed")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(suite.db)
	tableRepo := repositories.NewTableRepository(suite.db)
	menuRepo := repositories.NewMenuRepository(suite.db)
	orderRepo := repositories.NewOrderRepository(suite.db)
	paymentRepo := repositories.NewPaymentRepository(suite.db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
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
	userController := controllers.NewUserController(userService)
	orderManagementController := controllers.NewOrderManagementController(orderService)
	paymentManagementController := controllers.NewPaymentManagementController(paymentService)

	// Setup router with all controllers
	suite.router = setupTestRouter(cfg, authController, tableController, menuController, orderController, paymentController, kitchenController, userController, orderManagementController, paymentManagementController)

	// Create test users and get tokens
	suite.createTestUsers(authService)
	
	fmt.Println("SetupSuite completed successfully")
}

// Test method to verify suite is working
func (suite *APITestSuite) TestSetupWorking() {
	suite.T().Log("TestSetupWorking called")
	suite.NotNil(suite.db, "Database should be initialized")
	suite.NotNil(suite.router, "Router should be initialized")
}

// TearDownSuite runs after all tests
func (suite *APITestSuite) TearDownSuite() {
	// Clean up test database
	suite.dropTestTables()
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

func setupTestDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models for testing
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

func (suite *APITestSuite) createTestUsers(authService *services.AuthService) {
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Create admin user
	adminUser := &repositories.User{
		Username: "admin",
		Email:    "admin@test.com",
		Password: string(hashedPassword),
		Role:     repositories.RoleAdmin,
	}
	suite.db.Create(adminUser)

	// Create staff user
	staffUser := &repositories.User{
		Username: "staff",
		Email:    "staff@test.com",
		Password: string(hashedPassword),
		Role:     repositories.RoleStaff,
	}
	suite.db.Create(staffUser)

	// Create cashier user
	cashierUser := &repositories.User{
		Username: "cashier",
		Email:    "cashier@test.com",
		Password: string(hashedPassword),
		Role:     repositories.RoleCashier,
	}
	suite.db.Create(cashierUser)

	// Create customer user
	customerUser := &repositories.User{
		Username: "customer",
		Email:    "customer@test.com",
		Password: string(hashedPassword),
		Role:     repositories.RoleCustomer,
	}
	suite.db.Create(customerUser)

	// Generate tokens using login
	adminLogin := &services.LoginRequest{Username: "admin", Password: password}
	adminAuth, _ := authService.Login(adminLogin)
	suite.adminToken = adminAuth.AccessToken

	staffLogin := &services.LoginRequest{Username: "staff", Password: password}
	staffAuth, _ := authService.Login(staffLogin)
	suite.staffToken = staffAuth.AccessToken

	cashierLogin := &services.LoginRequest{Username: "cashier", Password: password}
	cashierAuth, _ := authService.Login(cashierLogin)
	suite.cashierToken = cashierAuth.AccessToken

	customerLogin := &services.LoginRequest{Username: "customer", Password: password}
	customerAuth, _ := authService.Login(customerLogin)
	suite.customerToken = customerAuth.AccessToken
}

func (suite *APITestSuite) dropTestTables() {
	suite.db.Migrator().DropTable(
		&repositories.Payment{},
		&repositories.OrderItem{},
		&repositories.Order{},
		&repositories.MenuItem{},
		&repositories.MenuCategory{},
		&repositories.Table{},
		&repositories.User{},
	)
}

// Helper methods
func (suite *APITestSuite) makeRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *APITestSuite) parseResponse(w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	suite.Require().NoError(err)
}

// Run the test suite
func TestAPITestSuite(t *testing.T) {
	t.Log("Starting TestAPITestSuite")
	s := new(APITestSuite)
	t.Log("Created suite instance")
	suite.Run(t, s)
	t.Log("Finished TestAPITestSuite")
}

// setupTestRouter creates a test router with all endpoints
func setupTestRouter(cfg *config.Config, authController *controllers.AuthController, tableController *controllers.TableController, menuController *controllers.MenuController, orderController *controllers.OrderController, paymentController *controllers.PaymentController, kitchenController *controllers.KitchenController, userController *controllers.UserController, orderManagementController *controllers.OrderManagementController, paymentManagementController *controllers.PaymentManagementController) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	
	// Use minimal middleware for testing
	router.Use(gin.Recovery())

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

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
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
				tables.PATCH("/:id/availability", tableController.UpdateTableAvailability)
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
		}

		// Staff routes (staff and admin access)
		staff := api.Group("/staff")
		staff.Use(middleware.AuthMiddleware())
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
		cashier.Use(middleware.AuthMiddleware())
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

	return router
}
