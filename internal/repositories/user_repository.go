package repositories

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserFilters struct {
	Role   string
	Search string
	Status string
}

type UserStatistics struct {
	TotalUsers    int64
	ActiveUsers   int64
	InactiveUsers int64
	UsersByRole   map[string]int64
	RecentUsers   []*User
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
		if strings.Contains(err.Error(), "phone") && strings.Contains(err.Error(), "does not exist") {
			return false, nil
		}
		return false, err
	}
	return count > 0, err
}

func (r *UserRepository) GetAll(limit, offset int) ([]User, error) {
	var users []User
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error
	for i := range users {
		users[i].Password = ""
	}
	return users, err
}

func (r *UserRepository) GetAllWithFilters(page, limit int, filters UserFilters) ([]*User, int64, error) {
	var users []*User
	var total int64
	offset := (page - 1) * limit
	query := r.db.Model(&User{}).Select("id, username, email, name, phone, role, is_active, created_at, updated_at")
	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}
	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(username) LIKE ? OR LOWER(email) LIKE ?", searchTerm, searchTerm, searchTerm)
	}
	switch filters.Status {
	case "active":
		query = query.Where("is_active = ?", true)
	case "inactive":
		query = query.Where("is_active = ?", false)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) GetUserStatistics() (*UserStatistics, error) {
	stats := &UserStatistics{UsersByRole: make(map[string]int64)}
	if err := r.db.Model(&User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}
	stats.InactiveUsers = stats.TotalUsers - stats.ActiveUsers
	type RoleCount struct {
		Role  string
		Count int64
	}
	var roleCounts []RoleCount
	if err := r.db.Model(&User{}).Select("role, COUNT(*) as count").Group("role").Scan(&roleCounts).Error; err != nil {
		return nil, err
	}
	for _, rc := range roleCounts {
		stats.UsersByRole[rc.Role] = rc.Count
	}
	var recentUsers []*User
	if err := r.db.Select("id, username, email, name, phone, role, is_active, created_at, updated_at").Order("created_at DESC").Limit(10).Find(&recentUsers).Error; err != nil {
		return nil, err
	}
	stats.RecentUsers = recentUsers
	return stats, nil
}

func (r *UserRepository) UserExistsByUsername(username string) (bool, error) {
	return r.IsUsernameExists(username)
}

func (r *UserRepository) UserExistsByEmail(email string) (bool, error) {
	return r.IsEmailExists(email)
}

func (r *UserRepository) UserExistsByPhone(phone string) (bool, error) {
	return r.IsPhoneExists(phone)
}
