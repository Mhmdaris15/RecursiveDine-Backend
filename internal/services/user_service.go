package services

import (
	"errors"
	"fmt"
	"strings"

	"recursiveDine/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// BulkUpdateResult represents the result of bulk operations
type BulkUpdateResult struct {
	UpdatedCount int      `json:"updated_count"`
	FailedCount  int      `json:"failed_count"`
	FailedIDs    []uint   `json:"failed_ids,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// GetAllUsers retrieves users with pagination (backward compatibility)
func (s *UserService) GetAllUsers(page, limit int) ([]repositories.User, error) {
	offset := (page - 1) * limit
	return s.userRepo.GetAll(limit, offset)
}

// GetAllUsersWithFilters retrieves users with pagination and filtering
func (s *UserService) GetAllUsersWithFilters(page, limit int, filters repositories.UserFilters) ([]*repositories.User, int64, error) {
	return s.userRepo.GetAllWithFilters(page, limit, filters)
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*repositories.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	
	// Remove password from response
	user.Password = ""
	return user, nil
}

// CreateUser creates a new user with proper validation
func (s *UserService) CreateUser(user *repositories.User) (*repositories.User, error) {
	// Validate required fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	// Validate role
	validRoles := []repositories.UserRole{
		repositories.RoleAdmin, 
		repositories.RoleCashier, 
		repositories.RoleStaff, 
		repositories.RoleCustomer,
	}
	if user.Role != "" && !containsRole(validRoles, user.Role) {
		return nil, errors.New("invalid role")
	}

	// Check if username already exists
	if exists, err := s.userRepo.IsUsernameExists(user.Username); err != nil {
		return nil, errors.New("failed to check username")
	} else if exists {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if exists, err := s.userRepo.IsEmailExists(user.Email); err != nil {
		return nil, errors.New("failed to check email")
	} else if exists {
		return nil, errors.New("email already exists")
	}

	// Check if phone already exists
	if user.Phone != "" {
		if exists, err := s.userRepo.IsPhoneExists(user.Phone); err != nil {
			return nil, errors.New("failed to check phone")
		} else if exists {
			return nil, errors.New("phone already exists")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	
	// Set default values
	if user.Role == "" {
		user.Role = repositories.RoleCustomer
	}
	user.IsActive = true

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// UpdateUser updates an existing user (backward compatibility)
func (s *UserService) UpdateUser(user *repositories.User) (*repositories.User, error) {
	// Check if user exists
	existing, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if username changed and if new username already exists
	if existing.Username != user.Username {
		if exists, err := s.userRepo.IsUsernameExists(user.Username); err != nil {
			return nil, errors.New("failed to check username")
		} else if exists {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email changed and if new email already exists
	if existing.Email != user.Email {
		if exists, err := s.userRepo.IsEmailExists(user.Email); err != nil {
			return nil, errors.New("failed to check email")
		} else if exists {
			return nil, errors.New("email already exists")
		}
	}

	// Check if phone changed and if new phone already exists
	if existing.Phone != user.Phone {
		if exists, err := s.userRepo.IsPhoneExists(user.Phone); err != nil {
			return nil, errors.New("failed to check phone")
		} else if exists {
			return nil, errors.New("phone already exists")
		}
	}

	// If password is provided, hash it
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	} else {
		// Keep existing password
		user.Password = existing.Password
	}

	// Validate role if provided
	if user.Role != "" {
		validRoles := []repositories.UserRole{
			repositories.RoleAdmin, 
			repositories.RoleCashier, 
			repositories.RoleStaff, 
			repositories.RoleCustomer,
		}
		if !containsRole(validRoles, user.Role) {
			return nil, errors.New("invalid role")
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// UpdateUserByAdmin updates a user with admin-specific logic
func (s *UserService) UpdateUserByAdmin(userID uint, updates interface{}) (*repositories.User, error) {
	// Get existing user
	existing, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Handle different update types
	var updateMap map[string]interface{}
	
	switch v := updates.(type) {
	case map[string]interface{}:
		updateMap = v
	default:
		// Try to handle struct-like updates
		updateMap = make(map[string]interface{})
		// This would need reflection for a complete implementation
		// For now, we'll handle the common cases
	}

	// Create a copy for updates
	updatedUser := *existing

	// Apply updates
	if username, ok := updateMap["username"].(string); ok && username != "" {
		if username != existing.Username {
			if exists, err := s.userRepo.IsUsernameExists(username); err != nil {
				return nil, errors.New("failed to check username")
			} else if exists {
				return nil, errors.New("username already exists")
			}
		}
		updatedUser.Username = username
	}

	if email, ok := updateMap["email"].(string); ok && email != "" {
		if email != existing.Email {
			if exists, err := s.userRepo.IsEmailExists(email); err != nil {
				return nil, errors.New("failed to check email")
			} else if exists {
				return nil, errors.New("email already exists")
			}
		}
		updatedUser.Email = email
	}

	if name, ok := updateMap["name"].(string); ok && name != "" {
		updatedUser.Name = name
	}

	if phone, ok := updateMap["phone"].(string); ok && phone != "" {
		if phone != existing.Phone {
			if exists, err := s.userRepo.IsPhoneExists(phone); err != nil {
				return nil, errors.New("failed to check phone")
			} else if exists {
				return nil, errors.New("phone already exists")
			}
		}
		updatedUser.Phone = phone
	}

	if password, ok := updateMap["password"].(string); ok && password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		updatedUser.Password = string(hashedPassword)
	}

	if role, ok := updateMap["role"].(string); ok && role != "" {
		validRoles := []repositories.UserRole{
			repositories.RoleAdmin, 
			repositories.RoleCashier, 
			repositories.RoleStaff, 
			repositories.RoleCustomer,
		}
		userRole := repositories.UserRole(role)
		if !containsRole(validRoles, userRole) {
			return nil, errors.New("invalid role")
		}
		updatedUser.Role = userRole
	}

	if err := s.userRepo.Update(&updatedUser); err != nil {
		return nil, errors.New("failed to update user")
	}

	// Remove password from response
	updatedUser.Password = ""
	return &updatedUser, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

// UpdateUserStatus updates user activation status
func (s *UserService) UpdateUserStatus(id uint, isActive bool) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.IsActive = isActive
	return s.userRepo.Update(user)
}

// UpdateUserRole updates user role
func (s *UserService) UpdateUserRole(userID uint, role string) (*repositories.User, error) {
	// Validate role
	validRoles := []repositories.UserRole{
		repositories.RoleAdmin, 
		repositories.RoleCashier, 
		repositories.RoleStaff, 
		repositories.RoleCustomer,
	}
	userRole := repositories.UserRole(role)
	if !containsRole(validRoles, userRole) {
		return nil, errors.New("invalid role")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Role = userRole
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user role")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

// ResetUserPassword resets a user's password
func (s *UserService) ResetUserPassword(userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// GetUserStatistics returns user statistics
func (s *UserService) GetUserStatistics() (*repositories.UserStatistics, error) {
	stats, err := s.userRepo.GetUserStatistics()
	if err != nil {
		return nil, fmt.Errorf("failed to get user statistics: %v", err)
	}

	return stats, nil
}

// BulkUpdateUsers performs bulk updates on multiple users
func (s *UserService) BulkUpdateUsers(userIDs []uint, updates map[string]interface{}) (*BulkUpdateResult, error) {
	result := &BulkUpdateResult{
		FailedIDs: []uint{},
		Errors:    []string{},
	}

	// Validate updates
	allowedFields := map[string]bool{
		"is_active": true,
		"role":      true,
	}

	for field := range updates {
		if !allowedFields[field] {
			return nil, fmt.Errorf("field '%s' is not allowed for bulk updates", field)
		}
	}

	// Validate role if provided
	if role, ok := updates["role"].(string); ok {
		validRoles := []repositories.UserRole{
			repositories.RoleAdmin, 
			repositories.RoleCashier, 
			repositories.RoleStaff, 
			repositories.RoleCustomer,
		}
		userRole := repositories.UserRole(role)
		if !containsRole(validRoles, userRole) {
			return nil, errors.New("invalid role")
		}
	}

	// Process each user
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			result.Errors = append(result.Errors, fmt.Sprintf("User %d: not found", userID))
			continue
		}

		// Apply updates
		if isActive, ok := updates["is_active"].(bool); ok {
			user.IsActive = isActive
		}
		if role, ok := updates["role"].(string); ok {
			user.Role = repositories.UserRole(role)
		}

		err = s.userRepo.Update(user)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, userID)
			result.Errors = append(result.Errors, fmt.Sprintf("User %d: %v", userID, err))
		} else {
			result.UpdatedCount++
		}
	}

	return result, nil
}

// SearchUsers searches for users by various criteria
func (s *UserService) SearchUsers(query string, filters repositories.UserFilters, page, limit int) ([]*repositories.User, int64, error) {
	// Clean up search query
	query = strings.TrimSpace(query)
	
	// Combine query with filters
	if query != "" && filters.Search == "" {
		filters.Search = query
	}

	return s.userRepo.GetAllWithFilters(page, limit, filters)
}

// Helper function to check if a slice contains a UserRole
func containsRole(slice []repositories.UserRole, item repositories.UserRole) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
