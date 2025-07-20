package repositories

import (
	"errors"

	"gorm.io/gorm"
)

type TableRepository struct {
	db *gorm.DB
}

func NewTableRepository(db *gorm.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(table *Table) error {
	return r.db.Create(table).Error
}

func (r *TableRepository) GetByID(id uint) (*Table, error) {
	var table Table
	err := r.db.First(&table, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("table not found")
	}
	return &table, err
}

func (r *TableRepository) GetByQRCode(qrCode string) (*Table, error) {
	var table Table
	err := r.db.Where("qr_code = ?", qrCode).First(&table).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("table not found")
	}
	return &table, err
}

func (r *TableRepository) GetByNumber(number int) (*Table, error) {
	var table Table
	err := r.db.Where("number = ?", number).First(&table).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("table not found")
	}
	return &table, err
}

func (r *TableRepository) GetAllAvailable() ([]Table, error) {
	var tables []Table
	err := r.db.Where("is_available = ?", true).Find(&tables).Error
	return tables, err
}

func (r *TableRepository) Update(table *Table) error {
	return r.db.Save(table).Error
}

func (r *TableRepository) Delete(id uint) error {
	return r.db.Delete(&Table{}, id).Error
}

func (r *TableRepository) IsQRCodeExists(qrCode string) (bool, error) {
	var count int64
	err := r.db.Model(&Table{}).Where("qr_code = ?", qrCode).Count(&count).Error
	return count > 0, err
}

func (r *TableRepository) IsNumberExists(number int) (bool, error) {
	var count int64
	err := r.db.Model(&Table{}).Where("number = ?", number).Count(&count).Error
	return count > 0, err
}

func (r *TableRepository) GetAll(limit, offset int) ([]Table, error) {
	var tables []Table
	err := r.db.Order("number ASC").
		Limit(limit).
		Offset(offset).
		Find(&tables).Error
	return tables, err
}
