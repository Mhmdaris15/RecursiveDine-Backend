package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"
	"recursiveDine/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// UserResponse represents the response structure for user data
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUserRequest represents the request structure for creating users
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,min=10,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=customer staff cashier admin"`
}

// UpdateUserRequest represents the request structure for updating users
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Name     string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Phone    string `json:"phone,omitempty" binding:"omitempty,min=10,max=20"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
	Role     string `json:"role,omitempty" binding:"omitempty,oneof=customer staff cashier admin"`
}

// @Summary Get all users
// @Description Get all users with pagination and filtering (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Param role query string false "Filter by role"
// @Param search query string false "Search by name, username, or email"
// @Param status query string false "Filter by status: active, inactive, all (default: all)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	role := c.Query("role")
	search := c.Query("search")
	status := c.DefaultQuery("status", "all")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filters := repositories.UserFilters{
		Role:   role,
		Search: strings.TrimSpace(search),
		Status: status,
	}

	users, total, err := ctrl.userService.GetAllUsersWithFilters(page, limit, filters)
	if err != nil {
		utils.LogError("Failed to get users", err, map[string]interface{}{
			"page":    page,
			"limit":   limit,
			"filters": filters,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"filters":     filters,
	})
}

// @Summary Get user by ID
// @Description Get user details by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [get]
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(userID))
	if err != nil {
		utils.LogError("Failed to get user by ID", err, map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Create user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /admin/users [post]
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Convert request to user model
	user := &repositories.User{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
		Role:     repositories.UserRole(req.Role),
	}

	createdUser, err := ctrl.userService.CreateUser(user)
	if err != nil {
		utils.LogError("Failed to create user", err, map[string]interface{}{
			"username": req.Username,
			"email":    req.Email,
			"role":     req.Role,
		})
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("User created successfully", map[string]interface{}{
		"user_id":  createdUser.ID,
		"username": createdUser.Username,
		"role":     createdUser.Role,
	})

	c.JSON(http.StatusCreated, createdUser)
}

// @Summary Update user
// @Description Update user details (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "User data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [put]
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	updatedUser, err := ctrl.userService.UpdateUserByAdmin(uint(userID), req)
	if err != nil {
		utils.LogError("Failed to update user", err, map[string]interface{}{
			"user_id": userID,
			"request": req,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("User updated successfully", map[string]interface{}{
		"user_id":  updatedUser.ID,
		"username": updatedUser.Username,
	})

	c.JSON(http.StatusOK, updatedUser)
}

// @Summary Delete user
// @Description Soft delete user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prevent users from deleting themselves
	currentUserID, exists := c.Get("user_id")
	if exists && currentUserID.(uint) == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	if err := ctrl.userService.DeleteUser(uint(userID)); err != nil {
		utils.LogError("Failed to delete user", err, map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("User deleted successfully", map[string]interface{}{
		"user_id": userID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// @Summary Activate/Deactivate user
// @Description Activate or deactivate user account (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body map[string]bool true "Active status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id}/status [patch]
func (ctrl *UserController) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prevent users from deactivating themselves
	currentUserID, exists := c.Get("user_id")
	if exists && currentUserID.(uint) == uint(userID) && !req.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot deactivate your own account"})
		return
	}

	if err := ctrl.userService.UpdateUserStatus(uint(userID), req.IsActive); err != nil {
		utils.LogError("Failed to update user status", err, map[string]interface{}{
			"user_id":   userID,
			"is_active": req.IsActive,
		})
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	status := "deactivated"
	if req.IsActive {
		status = "activated"
	}

	utils.LogInfo("User status updated", map[string]interface{}{
		"user_id": userID,
		"status":  status,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User " + status + " successfully"})
}

// @Summary Update user role
// @Description Update user role (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body map[string]string true "Role data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id}/role [patch]
func (ctrl *UserController) UpdateUserRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=customer staff cashier admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role: " + err.Error()})
		return
	}

	// Prevent users from changing their own role
	currentUserID, exists := c.Get("user_id")
	if exists && currentUserID.(uint) == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change your own role"})
		return
	}

	updatedUser, err := ctrl.userService.UpdateUserRole(uint(userID), req.Role)
	if err != nil {
		utils.LogError("Failed to update user role", err, map[string]interface{}{
			"user_id": userID,
			"role":    req.Role,
		})
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("User role updated", map[string]interface{}{
		"user_id":  userID,
		"new_role": req.Role,
	})

	c.JSON(http.StatusOK, updatedUser)
}

// @Summary Reset user password
// @Description Reset user password (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body map[string]string true "New password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/users/{id}/password [patch]
func (ctrl *UserController) ResetUserPassword(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password: " + err.Error()})
		return
	}

	if err := ctrl.userService.ResetUserPassword(uint(userID), req.Password); err != nil {
		utils.LogError("Failed to reset user password", err, map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("User password reset", map[string]interface{}{
		"user_id": userID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// @Summary Get user statistics
// @Description Get user statistics (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/statistics [get]
func (ctrl *UserController) GetUserStatistics(c *gin.Context) {
	stats, err := ctrl.userService.GetUserStatistics()
	if err != nil {
		utils.LogError("Failed to get user statistics", err, nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Bulk update users
// @Description Bulk update multiple users (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Bulk update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/users/bulk [patch]
func (ctrl *UserController) BulkUpdateUsers(c *gin.Context) {
	var req struct {
		UserIDs []uint                 `json:"user_ids" binding:"required"`
		Updates map[string]interface{} `json:"updates" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user IDs provided"})
		return
	}

	result, err := ctrl.userService.BulkUpdateUsers(req.UserIDs, req.Updates)
	if err != nil {
		utils.LogError("Failed to bulk update users", err, map[string]interface{}{
			"user_ids": req.UserIDs,
			"updates":  req.Updates,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("Bulk update completed", map[string]interface{}{
		"updated_count": result.UpdatedCount,
		"failed_count":  result.FailedCount,
	})

	c.JSON(http.StatusOK, result)
}
