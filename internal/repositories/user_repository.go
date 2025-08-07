package repositories

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *UserRepository) IsUsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) IsPhoneExists(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		// If phone column doesn't exist, return false (phone is available)
		if err.Error() == "column \"phone\" does not exist" || 
		   err.Error() == "pq: column \"phone\" does not exist" {
			return false, nil
		}
		return false, err
	}
	return count > 0, err
}

func (r *UserRepository) GetAll(limit, offset int) ([]User, error) {
	var users []User
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	
	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}
	
	return users, err
}
