package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"recursiveDine/internal/config"
	"recursiveDine/internal/repositories"
	"recursiveDine/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required,min=10,max=20"`
	Role     string `json:"role,omitempty"` // Optional role field for admin/staff registration
}

type AuthResponse struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	TokenType    string              `json:"token_type"`
	ExpiresIn    int                 `json:"expires_in"`
	User         *repositories.User  `json:"user"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user *repositories.User
	var err error

	// Support login with either username or email
	if req.Username != "" {
		user, err = s.userRepo.GetByUsername(req.Username)
	} else if req.Email != "" {
		user, err = s.userRepo.GetByEmail(req.Email)
	} else {
		return nil, errors.New("username or email is required")
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.config.JWTExpirationHours * 3600,
		User:         user,
	}, nil
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	if exists, err := s.userRepo.IsEmailExists(req.Email); err != nil {
		return nil, errors.New("failed to check email")
	} else if exists {
		return nil, errors.New("email already exists")
	}

	// Check if phone already exists
	if exists, err := s.userRepo.IsPhoneExists(req.Phone); err != nil {
		return nil, errors.New("failed to check phone")
	} else if exists {
		return nil, errors.New("phone already exists")
	}

	// Generate username from email (part before @)
	username := req.Email[:strings.Index(req.Email, "@")]
	
	// Make username unique if it already exists
	originalUsername := username
	counter := 1
	for {
		if exists, err := s.userRepo.IsUsernameExists(username); err != nil {
			return nil, errors.New("failed to check username")
		} else if !exists {
			break
		}
		username = fmt.Sprintf("%s%d", originalUsername, counter)
		counter++
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Determine user role
	userRole := repositories.RoleCustomer // Default role
	if req.Role != "" {
		// Validate role
		switch req.Role {
		case "admin":
			userRole = repositories.RoleAdmin
		case "staff":
			userRole = repositories.RoleStaff
		case "cashier":
			userRole = repositories.RoleCashier
		case "customer":
			userRole = repositories.RoleCustomer
		default:
			return nil, errors.New("invalid role specified")
		}
	}

	// Create user
	user := &repositories.User{
		Name:     req.Name,
		Username: username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Role:     userRole,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.config.JWTExpirationHours * 3600,
		User:         user,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Remove password from response
	user.Password = ""

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.config.JWTExpirationHours * 3600,
		User:         user,
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*repositories.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *AuthService) generateAccessToken(user *repositories.User) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.JWTExpirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   utils.UintToString(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *repositories.User) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.JWTRefreshHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   utils.UintToString(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}
