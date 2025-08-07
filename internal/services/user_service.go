package services

import (
	"errors"

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

func (s *UserService) GetAllUsers(page, limit int) ([]repositories.User, error) {
	offset := (page - 1) * limit
	return s.userRepo.GetAll(limit, offset)
}

func (s *UserService) GetUserByID(id uint) (*repositories.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *UserService) CreateUser(user *repositories.User) (*repositories.User, error) {
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
	if exists, err := s.userRepo.IsPhoneExists(user.Phone); err != nil {
		return nil, errors.New("failed to check phone")
	} else if exists {
		return nil, errors.New("phone already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.IsActive = true

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

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

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *UserService) UpdateUserStatus(id uint, isActive bool) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.IsActive = isActive
	return s.userRepo.Update(user)
}
